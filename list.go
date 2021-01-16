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
	"fmt"
	"path/filepath"
	"strings"
)

func list(context ExecutionContext) error {

	pathArg := "."

	if len(context.Args) > 0 {
		pathArg = context.Args[0]
	}

	showFullPath := *context.Options["fullpath"].(*bool)
	recursive := *context.Options["recursive"].(*bool)

	matches, err := resolveToGroupedMatches(pathArg, recursive)

	if err != nil {
		return err
	}

	rootPathLastIndex := len(filepath.Dir(commonBasePath(sortedKeys(matches)))) + 1
	firstPath := true

	for _, path := range sortedKeys(matches) {
		if showFullPath {
			for _, m := range matches[path] {
				fmt.Fprintln(out, m)
			}
		} else {
			if len(matches) > 1 {
				if !firstPath {
					fmt.Fprintln(out, "")
				} else {
					firstPath = false
				}
				fmt.Fprintln(out, path[rootPathLastIndex:]+":")
			}
			for _, m := range matches[path] {
				if isDir(m) {
					fmt.Fprintln(out, "+ "+filepath.Base(m))
				} else {
					fmt.Fprintln(out, "- "+filepath.Base(m)[:len(filepath.Base(m))-5])
				}
			}
		}
	}

	return nil

}

func commonBasePath(paths []string) string {
	commonPath := ""
	found := false
	for i := 1; !found && i < len(shortestPath(paths)); i++ {
		testPath := paths[0][:i]
		for _, path := range paths {
			if strings.Index(path, testPath) != 0 {
				found = true
			}
		}
		if !found {
			commonPath = testPath
		}
	}
	return commonPath
}

func shortestPath(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	shortest := paths[0]
	for _, path := range paths {
		if len(path) < len(shortest) {
			shortest = path
		}
	}
	return shortest
}
