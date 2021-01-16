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

func makeFolder(context ExecutionContext) error {

	if len(context.Args) < 1 {
		return errors.New("path required")
	}

	if !isValidPresetFolderName(context.Args[0]) {
		return errors.New("invalid folder name")
	}

	if isDir(context.Args[0]) {
		return errors.New("folder already exists")
	}

	target, _ := filepath.Abs(context.Args[0])

	err := os.MkdirAll(target, 0775)

	if err != nil {
		return errors.New("Failed to create folder: " + err.Error())
	}

	return nil

}
