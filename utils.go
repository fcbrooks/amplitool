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
	"os"
	"path/filepath"
)

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isSubDir(p1 string, p2 string) bool {
	if containsWildcards(p1) {
		return false
	}
	result := false
	exhausted := false
	for p2 != "" && !result && !exhausted {
		result = p1 == p2
		lastP2 := p2
		p2 = filepath.Dir(p2)
		exhausted = lastP2 == p2
	}
	return result
}

func isEmpty(folder string, recursive bool) bool {
	if isDir(folder) {
		if !recursive {
			matches, err := filepath.Glob(filepath.Join(folder, "*"))
			return err == nil && len(matches) == 0
		} else {
			result := true
			filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
				result = isDir(path)
				return nil
			})
			return result
		}
	}
	return false
}

func isV5Database(database *sql.DB) bool {
	_, checkVerErr := database.Prepare("select UserId, Product, Favorite, Description, Downloads, Keywords, Song, ChainA, ChainB, Band, Artist, ATInstrumentsType, ATPickupType, ATPickupPositions, ATSoundCharacter, ATGenre, SongStructureElement, Rating, MadeWith, ChainType, ATInstrument, OriginalFileName, FileFolder, Name from pXcPresets")
	return checkVerErr == nil
}
