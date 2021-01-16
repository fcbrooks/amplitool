/*
Copyright (C) 2021 fcbrooks

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func copy(context ExecutionContext) error {

	if len(context.Args) < 2 {
		return errors.New("missing source and destination paths.")
	}

	recursive := *context.Options["recursive"].(*bool)

	matches, err := resolveToMatches(context.Args[0], recursive, !containsWildcards(context.Args[0]) || recursive)

	if err != nil {
		return err
	}

	source, _ := filepath.Abs(context.Args[0])
	target, _ := filepath.Abs(context.Args[1])

	if !isInPresetsFolder(target) {
		return errors.New("presets not found on path")
	}

	sourceProfile, _ := resolveToProfile(source)
	targetProfile, _ := resolveToProfile(target)

	if isDir(source) && isValidPresetFolderName(target) && isSubDir(source, target) {
		return errors.New("cannot copy folder to subfolder of itself")
	}

	if isDir(source) && len(matches) > 0 && !recursive {
		fmt.Fprintln(out, "cp: -r not specified; omitting directory")
		return nil
	}

	if sourceProfile != targetProfile {
		return errors.New("source and target are not in the same profile")
	}

	files := map[string]string{}

	for _, match := range matches {
		if isValidPresetFolderName(target) {
			os.MkdirAll(target, 0775)
			if recursive {
				files[match] = filepath.Join(target, match[len(filepath.Dir(source)):])
			} else {
				if isDir(match) {
					fmt.Fprintln(out, "cp: -r not specified; omitting directory")
				} else {
					files[match] = filepath.Join(target, filepath.Base(match))
				}
			}
		} else {
			files[match] = target
		}
	}

	removeExistingDBRecord, err := context.Database.Prepare("delete from pXcPresets where OriginalFileName = ?")

	statement, err := context.Database.Prepare("insert into pXcPresets (UserId, Product, OriginalFileName, FileFolder, Favorite, Date, Name, Description, Downloads, Keywords, Song, ChainA, ChainB, Band, Artist, ATInstrumentsType, ATPickupType, ATPickupPositions, ATSoundCharacter, ATGenre, SongStructureElement, Rating, MadeWith, ChainType, tstamp, ATInstrument) select UserId, Product, ?, ?, Favorite, Date, ?, Description, Downloads, Keywords, Song, ChainA, ChainB, Band, Artist, ATInstrumentsType, ATPickupType, ATPickupPositions, ATSoundCharacter, ATGenre, SongStructureElement, Rating, MadeWith, ChainType, tstamp, ATInstrument from pXcPresets where OriginalFileName = ?")

	if err != nil {
		return errors.New("Failed preparing statement: " + err.Error())
	}

	for source, target := range files {

		if source == target {
			continue
		}

		overwritingExisting := isFile(target)

		err = os.MkdirAll(filepath.Dir(target), 0775)

		if err != nil {
			rollbackCopy(files)
			return errors.New("Failed to created directories :" + err.Error())
		}

		data, err := ioutil.ReadFile(source)

		if err != nil {
			rollbackCopy(files)
			return errors.New("Could not read source file: " + err.Error())
		}

		err = ioutil.WriteFile(target, data, 0644)

		if err != nil {
			rollbackCopy(files)
			return errors.New("Could not write file: " + err.Error())
		}

		err = writeNewGuidToFile(target)

		if err != nil {
			rollbackCopy(files)
			return errors.New("Failed to generate new GUID for copied file: " + err.Error())
		}

		if overwritingExisting {
			_, err = removeExistingDBRecord.Exec(target)
		}

		_, err = statement.Exec(target, filepath.Dir(target), makePresetPath(filepath.Base(target)), source)

		if err != nil {
			statement.Close()
			rollbackCopy(files)
			return errors.New("Copy failed.  Failed to update database: " + err.Error())
		}

	}

	statement.Close()

	return nil
}

func rollbackCopy(files map[string]string) {
	for _, target := range files {
		_, err := os.Stat(target)
		if err == nil {
			os.Remove(target)
		}
	}
}
