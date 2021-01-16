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
	"path/filepath"
)

func reindex(context ExecutionContext) error {

	if len(context.Args) == 0 {
		return errors.New("arg must be the root of an Amplitube profile")
	}

	source, _ := filepath.Abs(context.Args[0])

	if !isProfileFolder(source) {
		return errors.New("arg must be the root of an Amplitube profile")
	}

	database := context.Database

	statement, err := database.Prepare("update pXcPresets set OriginalFileName = ? || substr(OriginalFileName, instr(OriginalFileName, ?)), FileFolder = ? || substr(FileFolder, instr(FileFolder, ?))")

	if err != nil {
		return err
	}

	_, err = statement.Exec(source, filepath.Join(string(filepath.Separator), "Presets", string(filepath.Separator)), source, filepath.Join(string(filepath.Separator), "Presets", string(filepath.Separator)))

	if err != nil {
		return err
	}

	return nil

}
