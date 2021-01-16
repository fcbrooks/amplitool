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
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func importPresets(context ExecutionContext) error {

	sourcePath, _ := filepath.Abs(context.Args[0])
	targetPath, _ := filepath.Abs(context.Args[1])

	sourceProfile, err := resolveToProfile(sourcePath)

	if err != nil {
		return errors.New("both args must be the root of an Amplitube profile")
	}

	targetProfile, err := resolveToProfile(targetPath)

	if err != nil {
		return errors.New("both args must be the root of an Amplitube profile")
	}

	if sourceProfile == targetProfile {
		return errors.New("args cannot reference the same Amplitube profile")
	}

	matches, err := resolveToMatches(filepath.Join(sourceProfile, PresetsFolder), true, false)

	if err != nil {
		return err
	}

	importPath := filepath.Join(targetProfile, PresetsFolder, "Import")
	var counter int
	for info, _ := os.Stat(importPath); info != nil; info, _ = os.Stat(importPath) {
		counter++
		importPath = importPath + " - " + strconv.FormatInt(int64(counter), 10)
	}

	files := map[string]string{}

	for _, match := range matches {
		files[match] = filepath.Join(importPath, match[len(sourceProfile)+len(PresetsFolder)+1:])
	}

	sourceDatabase, err := openDatabase(filepath.Join(sourceProfile, "Presets.db"))

	if err != nil {
		return errors.New("failed to up source database")
	}

	var sourceStmt *sql.Stmt

	if !isV5Database(sourceDatabase) {
		sourceStmt, err = sourceDatabase.Prepare("select userid, product, 0 as favorite, description, downloads, keywords, song, chaina, chainb, NULL as band, NULL as artist, NULL as atinstrumentstype, NULL as atpickuptype, NULL as atpickuppositions, NULL as atsoundcharacter, NULL as atgenre, songstructureelement, rating, madewith, chaintype, NULL as atinstrument from pXcPresets where substr(OriginalFileName, instr(OriginalFileName, '" + string(filepath.Separator) + "Presets" + string(filepath.Separator) + "')) = substr(?, instr(?, '" + string(filepath.Separator) + "Presets" + string(filepath.Separator) + "'))")
	} else {
		sourceStmt, err = sourceDatabase.Prepare("select userid, product, favorite, description, downloads, keywords, song, chaina, chainb, band, artist, atinstrumentstype, atpickuptype, atpickuppositions, atsoundcharacter, atgenre, songstructureelement, rating, madewith, chaintype, atinstrument from pXcPresets where substr(OriginalFileName, instr(OriginalFileName, '" + string(filepath.Separator) + "Presets" + string(filepath.Separator) + "')) = substr(?, instr(?, '" + string(filepath.Separator) + "Presets" + string(filepath.Separator) + "'))")
	}

	if err != nil {
		return errors.New("Failed preparing statement: " + err.Error())
	}

	statement, err := context.Database.Prepare("insert into pXcPresets (UserId, Product, Favorite, Description, Downloads, Keywords, Song, ChainA, ChainB, Band, Artist, ATInstrumentsType, ATPickupType, ATPickupPositions, ATSoundCharacter, ATGenre, SongStructureElement, Rating, MadeWith, ChainType, ATInstrument, OriginalFileName, FileFolder, Name) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return errors.New("Failed preparing statement: " + err.Error())
	}

	for source, target := range files {

		if source == target {
			continue
		}

		if isValidPresetFolderName(target) {
			err = os.MkdirAll(target, 0775)
		} else {
			err = os.MkdirAll(filepath.Dir(target), 0775)
		}

		if err != nil {
			rollbackImport(importPath)
			return errors.New("Failed to created directories :" + err.Error())
		}

		if !isDir(target) {

			data, err := ioutil.ReadFile(source)

			if err != nil {
				rollbackImport(importPath)
				return errors.New("Could not read source file: " + err.Error())
			}

			err = ioutil.WriteFile(target, data, 0644)

			if err != nil {
				rollbackImport(importPath)
				return errors.New("Could not write file: " + err.Error())
			}

			err = writeNewGuidToFile(target)

			if err != nil {
				rollbackImport(importPath)
				return errors.New("Failed to generate new GUID for copied file: " + err.Error())
			}

			rs, err := sourceStmt.Query(source, source)

			if err != nil {
				sourceStmt.Close()
				sourceDatabase.Close()
				statement.Close()
				rollbackImport(importPath)
				return errors.New("Copy failed.  Failed to update database: " + err.Error())
			}

			inSourceDB := false
			var userid, product, favorite, description, downloads, keywords, song, chaina, chainb, band, artist, atinstrumentstype, atpickuptype, atpickuppositions, atsoundcharacter, atgenre, songstructureelement, rating, madewith, chaintype, atinstrument interface{}
			for rs.Next() {
				rs.Scan(&userid, &product, &favorite, &description, &downloads, &keywords, &song, &chaina, &chainb, &band, &artist, &atinstrumentstype, &atpickuptype, &atpickuppositions, &atsoundcharacter, &atgenre, &songstructureelement, &rating, &madewith, &chaintype, &atinstrument)
				inSourceDB = true
			}

			if inSourceDB {

				_, err = statement.Exec(userid, product, favorite, description, downloads, keywords, song, chaina, chainb, band, artist, atinstrumentstype, atpickuptype, atpickuppositions, atsoundcharacter, atgenre, songstructureelement, rating, madewith, chaintype, atinstrument, target, filepath.Dir(target), makePresetPath(filepath.Base(target)))

				if err != nil {
					sourceStmt.Close()
					sourceDatabase.Close()
					statement.Close()
					rollbackImport(importPath)
					return errors.New("Copy failed.  Failed to update database: " + err.Error())
				}

			} else {
				fmt.Fprintln(out, source+" imported without database record")
			}

		}

	}

	sourceStmt.Close()
	sourceDatabase.Close()
	statement.Close()

	return nil
}

func rollbackImport(importPath string) {
	os.RemoveAll(importPath)
}
