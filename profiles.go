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
)

func isProfileFolder(folder string) bool {
	info, err := os.Stat(filepath.Join(folder, PresetsFolder))
	if err == nil && info != nil {
		info, err = os.Stat(filepath.Join(folder, "Presets.db"))
		if err == nil && info != nil {
			return true
		}
	}
	return false
}

func resolveToProfile(startPath string) (string, error) {

	found := false
	exhausted := false
	currPath := startPath

	for !found && !exhausted {
		found = isProfileFolder(currPath)
		if !found {
			lastPath := currPath
			currPath = filepath.Dir(currPath)
			exhausted = lastPath == currPath
		}
	}

	if !found && exhausted {
		return "", errors.New(startPath + " is not part of Amplitube profile.")
	}

	return currPath, nil
}
