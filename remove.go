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

func remove(context ExecutionContext) error {

	if len(context.Args) < 1 {
		return errors.New("missing target path to remove")
	}

	source, _ := filepath.Abs(context.Args[0])
	recursive := *context.Options["recursive"].(*bool)

	if isDir(source) && !isEmpty(source, false) && !recursive {
		fmt.Fprintln(out, "-r not specified; omitting directory")
	}

	matches, err := resolveToMatches(context.Args[0], recursive, false)

	if err != nil {
		return err
	}

	var files []string
	var dirs []string

	for _, match := range matches {
		if info, _ := os.Stat(match); info.IsDir() {
			if !isEmpty(match, false) && !recursive {
				fmt.Fprintln(out, "-r not specified; omitting directory")
			} else {
				dirs = append(dirs, match)
			}
		} else {
			files = append(files, match)
		}
	}

	if isDir(source) && !containsWildcards(source) {
		dirs = append(dirs, source)
	}

	statement, err := context.Database.Prepare("delete from pXcPresets where OriginalFileName = ?")

	if err != nil {
		return errors.New("Failed preparing statement: " + err.Error())
	}

	for _, target := range files {

		info, err := os.Stat(target)

		if err != nil {
			rollbackRemove(files)
			return errors.New("Unable to stat " + target + ".  " + err.Error())
		}

		if !info.IsDir() {

			rollback := filepath.Join(filepath.Dir(target), "_"+filepath.Base(target))

			data, err := ioutil.ReadFile(target)

			if err != nil {
				rollbackRemove(files)
				return errors.New("Failed to backup file before delete: " + err.Error())
			}

			err = ioutil.WriteFile(rollback, data, 0644)

			if err != nil {
				rollbackRemove(files)
				return errors.New("Failed to backup file before delete: " + err.Error())
			}

			err = os.Remove(target)

			if err != nil {
				rollbackRemove(files)
				return errors.New("Failed to remove file: " + err.Error())
			}

			_, err = statement.Exec(target)

			if err != nil {
				statement.Close()
				rollbackRemove(files)
				return errors.New("Failed to remove file: " + err.Error())
			}
		}
	}

	for _, target := range files {
		rollback := filepath.Join(filepath.Dir(target), "_"+filepath.Base(target))
		os.Remove(rollback)
	}

	for _, target := range dirs {
		os.RemoveAll(target)
	}

	return nil

}

func rollbackRemove(files []string) {
	for _, target := range files {
		rollback := filepath.Join(filepath.Dir(target), "_"+filepath.Base(target))
		data, _ := ioutil.ReadFile(rollback)
		ioutil.WriteFile(target, data, 0644)
		os.Remove(rollback)
	}
}
