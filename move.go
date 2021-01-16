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
	"os"
	"path/filepath"
	"strings"
)

func move(context ExecutionContext) error {

	if len(context.Args) < 2 {
		return errors.New("missing source and destination paths.")
	}

	recursive := isValidPresetFolderName(context.Args[0]) && isValidPresetFolderName(context.Args[1])

	matches, err := resolveToMatches(context.Args[0], recursive, true)

	if err != nil {
		return err
	}

	source, _ := filepath.Abs(context.Args[0])
	target, _ := filepath.Abs(context.Args[1])

	if !isInPresetsFolder(target) {
		return errors.New("presets not found on path")
	}

	if isDir(source) && isValidPresetFolderName(target) && isSubDir(source, target) {
		return errors.New("cannot move folder to subfolder of itself")
	}

	sourceProfile, _ := resolveToProfile(source)
	targetProfile, _ := resolveToProfile(target)

	if sourceProfile != targetProfile {
		return errors.New("source and target are not in the same profile")
	}

	files := map[string]string{}

	for _, match := range matches {
		if isValidPresetFolderName(target) {
			os.MkdirAll(target, 0775)
			if recursive {
				if isDir(target) {
					if isDir(filepath.Join(target, filepath.Base(source))) {
						if isEmpty(filepath.Join(target, filepath.Base(source)), false) {
							files[match] = filepath.Join(target, match[len(filepath.Dir(source)):])
						} else {
							return errors.New("cannot move: Directory is not empty")
						}
					} else {
						files[match] = filepath.Join(target, match[len(filepath.Dir(source)):])
					}
				} else {
					files[match] = filepath.Join(target, match[len(source):])
				}
			} else {
				files[match] = filepath.Join(target, filepath.Base(match))
			}
		} else {
			files[match] = target
		}
	}

	statement, err := context.Database.Prepare("update pXcPresets set OriginalFileName = ?, FileFolder = ?, Name = ? where OriginalFileName = ?")

	if err != nil {
		return errors.New("Failed preparing statement: " + err.Error())
	}

	for source, target := range files {

		_, err := os.Stat(source)

		if err != nil {
			return errors.New("File not found: " + source + ".  " + err.Error())
		}

		err = os.MkdirAll(filepath.Dir(target), 0775)

		if err != nil {
			return errors.New("Failed to created directories :" + err.Error())
		}

		err = os.Rename(source, target)

		if err != nil {
			return errors.New("Failed to move preset with error: " + err.Error())
		}

		_, err = statement.Exec(target, filepath.Dir(target), makePresetPath(filepath.Base(target)), source)

		if err != nil {
			statement.Close()
			rollbackMove(files)
			return errors.New("Failed to update database records: " + err.Error())
		}

	}

	if isDir(source) {
		if isEmpty(source, true) {
			os.RemoveAll(source)
		} else {
			statement.Close()
			rollbackMove(files)
			return errors.New("move failed while removing empty folders")
		}
	} else {
		if containsWildcards(source) {
			cleanUp, _ := resolveToMatches(source, recursive, false)
			for _, file := range cleanUp {
				if file != target && strings.Index(file, target) != 0 {
					if !isDir(file) || isEmpty(file, true) {
						os.RemoveAll(file)
					} else {
						statement.Close()
						rollbackMove(files)
						return errors.New("move failed while removing empty folders")
					}
				}
			}
		}
	}

	statement.Close()

	return nil

}

func rollbackMove(files map[string]string) {
	for source, target := range files {
		_, err := os.Stat(target)
		if err == nil {
			os.Rename(target, source)
		}
	}
}
