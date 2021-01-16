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
	"sort"
	"strings"
)

func resolveToMatches(path string, recursive bool, excludeFolders bool) ([]string, error) {

	absPath, _ := filepath.Abs(path)

	chkPath := absPath

	for containsWildcards(filepath.Base(chkPath)) {
		chkPath = filepath.Dir(chkPath)
	}

	if info, _ := os.Stat(chkPath); info == nil {
		return []string{}, errors.New("path not found")
	}

	if !isInPresetsFolder(absPath) {
		return []string{}, errors.New("presets not found on path")
	}

	var matches []string

	if isValidPresetName(absPath) {
		matches = append(matches, absPath)
	} else {
		if recursive {
			if containsWildcards(absPath) {
				wcMatches, _ := filepath.Glob(absPath)
				for _, m := range wcMatches {
					if info, _ := os.Stat(m); !info.IsDir() {
						matches = append(matches, m)
					}
				}
				for _, wcPath := range wcMatches {
					if info, err := os.Stat(wcPath); err == nil && info.IsDir() {
						filepath.Walk(wcPath, func(walkedPath string, info os.FileInfo, err error) error {
							if absPath != walkedPath && !excludeFolders || !info.IsDir() {
								matches = append(matches, walkedPath)
							}
							return nil
						})
					}
				}
			} else {
				filepath.Walk(absPath, func(walkedPath string, info os.FileInfo, err error) error {
					if absPath != walkedPath && !excludeFolders || !info.IsDir() {
						matches = append(matches, walkedPath)
					}
					return nil
				})
			}
		} else {
			var possibles []string
			if containsWildcards(absPath) {
				possibles, _ = filepath.Glob(absPath)
			} else {
				possibles, _ = filepath.Glob(filepath.Join(absPath, "*"))
			}
			for _, possible := range possibles {
				if info, err := os.Stat(possible); err == nil && (!info.IsDir() || (info.IsDir() && !excludeFolders)) {
					matches = append(matches, possible)
				}
			}
		}
	}

	return matches, nil

}

func resolveToGroupedMatches(path string, recursive bool) (map[string][]string, error) {
	matches, err := resolveToMatches(path, recursive, false)
	return groupMatchesByFolder(matches), err
}

func groupMatchesByFolder(matches []string) map[string][]string {

	groupedMatches := map[string][]string{}

	for _, match := range matches {
		groupedMatches[filepath.Dir(match)] = append(groupedMatches[filepath.Dir(match)], match)
	}

	for k, v := range groupedMatches {
		groupedMatches[k] = sortMatches(v)
	}

	return groupedMatches

}

func sortMatches(matches []string) []string {
	sort.SliceStable(matches, func(i, j int) bool {
		if filepath.Dir(matches[i]) == filepath.Dir(matches[j]) && isDir(matches[i]) && !isDir(matches[j]) {
			return true
		}
		return matches[i] < matches[j]
	})
	return matches
}

func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func containsWildcards(path string) bool {
	return strings.Contains(path, "*")
}
