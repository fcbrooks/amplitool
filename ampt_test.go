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
	"bytes"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

const TestDataRoot = "testdata"
const EmptyPlaceholder = ".empty"

func TestAll(t *testing.T) {

	for _, tc := range []struct {
		Name             string
		Command          string
		Args             []string
		Expected         string
		ExpectedError    string
		ExpectExists     []string
		ExpectNotExist   []string
		ExpectDBExists   []string
		ExpectDBNotExist []string
		CustomSetup      func(workingDirs []string)
		CustomAssertion  func(workingDir string) error
		WorkDirCount     int
		Verbose          bool
	}{
		{
			Name:          "No path provided",
			Command:       "ls",
			ExpectedError: "presets not found on path",
		},
		{
			Name:    "Path is outside of profile",
			Command: "ls",
			Args: []string{
				".",
			},
			ExpectedError: "presets not found on path",
		},
		{
			Name:    "Path is top level of profile",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot),
			},
			ExpectedError: "presets not found on path",
		},
		{
			Name:    "List folders",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder),
			},
			Expected: `+ Amps
+ Amps2`,
		},
		{
			Name:    "List preset",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default.at5p"),
			},
			Expected: "- Default",
		},
		{
			Name:    "List preset full path",
			Command: "ls",
			Args: []string{
				"-f",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default.at5p"),
			},
			CustomAssertion: func(workingDir string) error {
				if !strings.Contains(out.(*bytes.Buffer).String(), filepath.Join(workingDir, PresetsFolder, "Amps", "Default.at5p")) {
					t.Errorf("wanted '%s'; was '%s'", filepath.Join(workingDir, PresetsFolder, "Amps", "Default.at5p"), out.(*bytes.Buffer).String())
				}
				return nil
			},
		},
		{
			Name:    "Preset does not exist",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Test.at5p"),
			},
			ExpectedError: "path not found",
		},
		{
			Name:    "Preset folder does not exist",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Test"),
			},
			ExpectedError: "path not found",
		},
		{
			Name:    "Presets after folder in list",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps"),
			},
			Expected: `+ Amplitube
+ Empty
+ THD
- Default`,
		},
		{
			Name:    "List preset full relative path recursive",
			Command: "ls",
			Args: []string{
				"-f",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps"),
			},
			CustomAssertion: func(workingDir string) error {
				if !strings.Contains(out.(*bytes.Buffer).String(), filepath.Join(workingDir, PresetsFolder, "Amps", "Default.at5p")) {
					t.Errorf("wanted '%s'; was '%s'", filepath.Join(workingDir, PresetsFolder, "Amps", "Default.at5p"), out.(*bytes.Buffer).String())
				}
				return nil
			},
		},
		{
			Name:    "Recursive list",
			Command: "ls",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
			},
			Expected: `Amplitube:
+ Metal
+ SVX
- American Tube Clean 1

Amplitube\Metal:
- Metal Clean T

Amplitube\SVX:
- SVX-4B`,
		},
		{
			Name:    "List preset full path recursive",
			Command: "ls",
			Args: []string{
				"-f",
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
			},
			CustomAssertion: func(workingDir string) error {
				expected := fmt.Sprintf(`%s
%s
%s`, filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1.at5p"),
					filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T.at5p"),
					filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B.at5p"),
				)
				if !strings.Contains(out.(*bytes.Buffer).String(), expected) {
					t.Errorf("wanted '%s'; was '%s'", expected, out.(*bytes.Buffer).String())
				}
				return nil
			},
		},
		{
			Name:    "List wildcard",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "*"),
			},
			Expected: `+ Amplitube
+ Empty
+ THD
- Default`,
		},
		{
			Name:    "List double wildcard",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*", "*"),
			},
			Expected: `Metal:
- Metal Clean T

SVX:
- SVX-4B`,
		},
		{
			Name:    "List wildcard recursive",
			Command: "ls",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"),
			},
			Expected: `Amplitube:
+ Metal
+ SVX
- American Tube Clean 1

Amplitube\Metal:
- Metal Clean T

Amplitube\SVX:
- SVX-4B`,
		},
		{
			Name:    "List wildcard narrower",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"+PresetExtension),
			},
			Expected: "- American Tube Clean 1",
		},
		{
			Name:    "List wildcard narrower recursive",
			Command: "ls",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"+PresetExtension),
			},
			Expected: "- American Tube Clean 1",
		},
		{
			Name:    "List wildcard prefix",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "A*"),
			},
			Expected: "- American Tube Clean 1",
		},
		{
			Name:    "List wildcard prefix",
			Command: "ls",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Met*"),
			},
			Expected: `+ Metal`,
		},
		{
			Name:    "List wildcard ",
			Command: "ls",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Met*"),
			},
			Expected: `Metal:
- Metal Clean T`,
		},
		{
			Name:          "Copy requires two parameters",
			Command:       "cp",
			ExpectedError: "missing source and destination paths",
		},
		{
			Name:    "Copy source not in profile",
			Command: "cp",
			Args: []string{
				".",
				filepath.Join(TestDataRoot, PresetsFolder, "Test"+PresetExtension),
			},
			ExpectedError: "presets not found on path",
		},
		{
			Name:    "Copy source preset root error",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder),
				filepath.Join(TestDataRoot, PresetsFolder, "Test"),
			},
			ExpectedError: "cannot copy folder to subfolder of itself",
		},
		{
			Name:    "Copy target in source error",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Test"),
			},
			ExpectedError: "cannot copy folder to subfolder of itself",
		},
		{
			Name:    "Copy target not in profile",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				".",
			},
			ExpectedError: "presets not found on path",
		},
		{
			Name:    "Copy one preset",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Test"+PresetExtension),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Test"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Test"+PresetExtension),
			},
		},
		{
			Name:    "Copy one preset over another",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension),
			},
			CustomAssertion: func(workingDir string) error {
				source := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				target := filepath.Join(workingDir, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension)
				sourceData, _ := ioutil.ReadFile(source)
				targetData, _ := ioutil.ReadFile(target)
				var sourceXml PresetXMLV5
				var targetXml PresetXMLV5
				xml.Unmarshal(sourceData, &sourceXml)
				xml.Unmarshal(targetData, &targetXml)
				if sourceXml.AmpA.Model != targetXml.AmpA.Model {
					return errors.New("source and target have different AmpA; expected same")
				}
				return nil
			},
		},
		{
			Name:    "Copy one preset to folder",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Default"+PresetExtension),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Default"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Default"+PresetExtension),
			},
		},
		{
			Name:    "Copied preset has new GUID",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Default"+PresetExtension),
			},
			CustomAssertion: func(workingDir string) error {
				source := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				target := filepath.Join(workingDir, PresetsFolder, "Amps", "Empty", "Default"+PresetExtension)
				sourceData, _ := ioutil.ReadFile(source)
				targetData, _ := ioutil.ReadFile(target)
				var sourceXml PresetXMLRootOnlyV5
				var targetXml PresetXMLRootOnlyV5
				xml.Unmarshal(sourceData, &sourceXml)
				xml.Unmarshal(targetData, &targetXml)
				if sourceXml.GUID == targetXml.GUID {
					return errors.New("source and target have same GUID; expected different")
				}
				return nil
			},
		},
		{
			Name:    "Copy folder with -r fails",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			Expected: "cp: -r not specified; omitting directory",
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube"),
			},
		},
		{
			Name:    "Recursively copy folder",
			Command: "cp",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Copy to target folder that does not exist",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Test"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Test", "Default"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Test", "Default"+PresetExtension),
			},
		},
		{
			Name:    "Wildcard copy to target folder that does not exist",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American*"),
				filepath.Join(TestDataRoot, PresetsFolder, "Test"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Test", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Test", "American Tube Clean 1"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Test", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Test", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Test", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Test", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Wildcard copy to target folder that matches",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American*"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American", "American Tube Clean 1"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Copy wildcard folder",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			Expected: "cp: -r not specified; omitting directory",
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Recursively copy folder with wildcard",
			Command: "cp",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Copy wildcard matching extension with recursive",
			Command: "cp",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Copy wildcard matching folder omitted",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "M*"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			Expected: "cp: -r not specified; omitting directory",
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Copy wildcard matching folder recursive",
			Command: "cp",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "M*"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Copy with Wildcard to folder in same parent that starts with same word",
			Command: "cp",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American*"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
		},
		{
			Name:          "Move requires two parameters",
			Command:       "mv",
			ExpectedError: "missing source and destination paths",
		},
		{
			Name:    "Move source not in profile",
			Command: "mv",
			Args: []string{
				".",
				filepath.Join(TestDataRoot, PresetsFolder, "Test"+PresetExtension),
			},
			ExpectedError: "presets not found on path",
		},
		{
			Name:    "Move source preset folder error",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder),
				filepath.Join(TestDataRoot, PresetsFolder, "Test"),
			},
			ExpectedError: "cannot move folder to subfolder of itself",
		},
		{
			Name:    "Move target in source folder error",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Test"),
			},
			ExpectedError: "cannot move folder to subfolder of itself",
		},
		{
			Name:    "Move target not in profile",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				".",
			},
			ExpectedError: "presets not found on path",
		},
		{
			Name:    "Move one preset",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Test"+PresetExtension),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Test"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Test"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
			},
		},
		{
			Name:    "Move one preset to folder",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Default"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Default"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
			},
		},
		{
			Name:    "Move folder",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Move to non-empty subfolder with same name fails",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps2"),
			},
			ExpectedError: "cannot move",
		},
		{
			Name:    "Move to empty subfolder with same name succeeds",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "THD"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps2"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps2", "THD", "BiValve"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "THD"),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps2", "THD", "BiValve"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension),
			},
		},
		{
			Name:    "Move with wildcard",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX"),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Move with wildcard matching extension",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Move with wildcard matching extension",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "M*"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal"),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "Metal", "Metal Clean T"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Empty", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Move with Wildcard to folder that starts with same word",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
		},
		{
			Name:    "Move with Wildcard to folder in same parent that starts with same word",
			Command: "mv",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American*"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
		},
		{
			Name:          "Remove requires one arg",
			Command:       "rm",
			ExpectedError: "missing target path to remove",
		},
		{
			Name:    "Remove target not in profile",
			Command: "rm",
			Args: []string{
				".",
			},
			ExpectedError: "presets not found on path",
		},
		{
			Name:    "Delete one preset",
			Command: "rm",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
		},
		{
			Name:    "Remove non-empty folder fails",
			Command: "rm",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal"),
			},
			Expected: "-r not specified; omitting directory",
		},
		{
			Name:    "Remove non-empty folder",
			Command: "rm",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal"),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal"),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
			},
		},
		{
			Name:    "Remove multi-level folder structure",
			Command: "rm",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps"),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Remove with wildcard",
			Command: "rm",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"),
			},
			Expected: "-r not specified; omitting directory",
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
		},
		{
			Name:    "Remove with wildcard matching extension",
			Command: "rm",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"+PresetExtension),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
		},
		{
			Name:    "Remove with wildcard matching extension",
			Command: "rm",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"+PresetExtension),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
			},
		},
		{
			Name:    "Remove with wildcard matching folder",
			Command: "rm",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "M*"),
			},
			Expected: "-r not specified; omitting directory",
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Remove with wildcard matching folder",
			Command: "rm",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "M*"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal"),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
			},
		},
		{
			Name:    "Remove with wildcard recursive",
			Command: "rm",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "*"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
			},
			ExpectNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX"),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:          "Path required to create folder",
			Command:       "mkdir",
			ExpectedError: "path required",
		},
		{
			Name:    "Invalid folder name",
			Command: "mkdir",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Test"+PresetExtension),
			},
			ExpectedError: "invalid folder name",
		},
		{
			Name:    "Create directory",
			Command: "mkdir",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Test"),
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Test"),
			},
		},
		{
			Name:    "Folder already exist",
			Command: "mkdir",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps"),
			},
			ExpectedError: "folder already exists",
		},
		{
			Name:          "Copy gear required three args",
			Command:       "cpg",
			ExpectedError: "copy gear requires a source preset, destination, and source gear name",
		},
		{
			Name:    "Copy gear source must be a file",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"AmpA",
			},
			ExpectedError: "first argument must be a preset file",
		},
		{
			Name:    "Copy gear source must be a preset",
			Command: "cpg",
			Args: []string{
				filepath.Join("README.md"),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"AmpA",
			},
			ExpectedError: "first argument must be a preset file",
		},
		{
			Name:    "Copy AmpA",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"AmpA",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.AmpA.Model != preset2.AmpA.Model {
					return errors.New("amps do not match; expected " + preset1.AmpA.Model + "; was " + preset2.AmpA.Model)
				}
				if !reflect.DeepEqual(preset1.AmpA.Amp, preset2.AmpA.Amp) {
					amp1, _ := xml.Marshal(preset1.AmpA.Amp)
					amp2, _ := xml.Marshal(preset1.AmpA.Amp)
					return errors.New("amp details do not match; expected " + string(selfClose(amp1)) + "; was " + string(selfClose(amp2)))
				}
				return nil
			},
		},
		{
			Name:    "Copy AmpB",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"AmpB",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.AmpB.Model != preset2.AmpB.Model {
					return errors.New("amps do not match; expected " + preset1.AmpB.Model + "; was " + preset2.AmpB.Model)
				}
				if !reflect.DeepEqual(preset1.AmpB.Amp, preset2.AmpB.Amp) {
					amp1, _ := xml.Marshal(preset1.AmpB.Amp)
					amp2, _ := xml.Marshal(preset1.AmpB.Amp)
					return errors.New("amp details do not match; expected " + string(selfClose(amp1)) + "; was " + string(selfClose(amp2)))
				}
				return nil
			},
		},
		{
			Name:    "Copy AmpC",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"AmpC",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.AmpC.Model != preset2.AmpC.Model {
					return errors.New("amps do not match; expected " + preset1.AmpC.Model + "; was " + preset2.AmpC.Model)
				}
				if !reflect.DeepEqual(preset1.AmpC.Amp, preset2.AmpC.Amp) {
					amp1, _ := xml.Marshal(preset1.AmpC.Amp)
					amp2, _ := xml.Marshal(preset1.AmpC.Amp)
					return errors.New("amp details do not match; expected " + string(selfClose(amp1)) + "; was " + string(selfClose(amp2)))
				}
				return nil
			},
		},
		{
			Name:    "Copy CabA",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"CabA",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.CabA.CabModel != preset2.CabA.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabA.CabModel + "; was " + preset2.CabA.CabModel)
				}
				if !reflect.DeepEqual(preset1.CabA.Cab, preset2.CabA.Cab) {
					cab1, _ := xml.Marshal(preset1.CabA.Cab)
					cab2, _ := xml.Marshal(preset2.CabA.Cab)
					return errors.New("cab details do not match; expected " + string(selfClose(cab1)) + "; was " + string(selfClose(cab2)))
				}
				return nil
			},
		},
		{
			Name:    "Copy CabB",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"CabB",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.CabB.CabModel != preset2.CabB.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabB.CabModel + "; was " + preset2.CabB.CabModel)
				}
				if !reflect.DeepEqual(preset1.CabB.Cab, preset2.CabB.Cab) {
					cab1, _ := xml.Marshal(preset1.CabB.Cab)
					cab2, _ := xml.Marshal(preset2.CabB.Cab)
					return errors.New("cab details do not match; expected " + string(selfClose(cab1)) + "; was " + string(selfClose(cab2)))
				}
				return nil
			},
		},
		{
			Name:    "Copy CabC",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"CabC",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.CabC.CabModel != preset2.CabC.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabC.CabModel + "; was " + preset2.CabC.CabModel)
				}
				if !reflect.DeepEqual(preset1.CabC.Cab, preset2.CabC.Cab) {
					cab1, _ := xml.Marshal(preset1.CabC.Cab)
					cab2, _ := xml.Marshal(preset2.CabC.Cab)
					return errors.New("cab details do not match; expected " + string(selfClose(cab1)) + "; was " + string(selfClose(cab2)))
				}
				return nil
			},
		},
		{
			Name:    "Copy one amp to different slot",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"AmpA",
				"AmpB",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.AmpA.Model != preset2.AmpB.Model {
					return errors.New("amps do not match; expected " + preset1.AmpA.Model + "; was " + preset2.AmpB.Model)
				}
				return nil
			},
		},
		{
			Name:    "Copy one cab to different slot",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"CabA",
				"CabB",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.CabA.CabModel != preset2.CabB.CabModel {
					return errors.New("amps do not match; expected " + preset1.CabA.CabModel + "; was " + preset2.CabB.CabModel)
				}
				return nil
			},
		},
		{
			Name:    "Copy Cab With Amp",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"AmpA",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.AmpA.Model != preset2.AmpA.Model {
					return errors.New("amps do not match; expected " + preset1.AmpA.Model + "; was " + preset2.AmpA.Model)
				}
				if preset1.CabA.CabModel != preset2.CabA.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabA.CabModel + "; was " + preset2.CabA.CabModel)
				}
				if !reflect.DeepEqual(preset1.AmpA.Amp, preset2.AmpA.Amp) {
					amp1, _ := xml.Marshal(preset1.AmpA.Amp)
					amp2, _ := xml.Marshal(preset1.AmpA.Amp)
					return errors.New("amp details do not match; expected " + string(selfClose(amp1)) + "; was " + string(selfClose(amp2)))
				}
				if !reflect.DeepEqual(preset1.CabA.CabModel, preset2.CabA.CabModel) {
					cab1, _ := xml.Marshal(preset1.CabA.CabModel)
					cab2, _ := xml.Marshal(preset2.CabA.CabModel)
					return errors.New("cab details do not match; expected " + string(selfClose(cab1)) + "; was " + string(selfClose(cab2)))
				}
				return nil
			},
		},
		{
			Name:    "Copy Cab With Amp To different slots",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"AmpA",
				"AmpB",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.AmpA.Model != preset2.AmpB.Model {
					return errors.New("amps do not match; expected " + preset1.AmpA.Model + "; was " + preset2.AmpA.Model)
				}
				if preset1.CabA.CabModel != preset2.CabB.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabA.CabModel + "; was " + preset2.CabB.CabModel)
				}
				if !reflect.DeepEqual(preset1.AmpA.Amp, preset2.AmpB.Amp) {
					amp1, _ := xml.Marshal(preset1.AmpA.Amp)
					amp2, _ := xml.Marshal(preset1.AmpB.Amp)
					return errors.New("amp details do not match; expected " + string(selfClose(amp1)) + "; was " + string(selfClose(amp2)))
				}
				if !reflect.DeepEqual(preset1.CabA.CabModel, preset2.CabB.CabModel) {
					cab1, _ := xml.Marshal(preset1.CabA.CabModel)
					cab2, _ := xml.Marshal(preset2.CabB.CabModel)
					return errors.New("cab details do not match; expected " + string(selfClose(cab1)) + "; was " + string(selfClose(cab2)))
				}
				return nil
			},
		},
		{
			Name:    "Don't Copy Cab With Amp",
			Command: "cpg",
			Args: []string{
				"-c",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"AmpA",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.AmpA.Model != preset2.AmpA.Model {
					return errors.New("amps do not match; expected " + preset1.AmpA.Model + "; was " + preset2.AmpA.Model)
				}
				if preset1.CabA.CabModel == preset2.CabA.CabModel {
					return errors.New("cabs match; both were " + preset1.CabA.CabModel)
				}
				return nil
			},
		},
		{
			Name:    "Overwrite Stomps same type",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"StompA1",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA1
				stomp2 := preset2.StompA1
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if stomp1.Stomp4 != stomp2.Stomp4 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp4 + "; was " + stomp2.Stomp4)
				}
				if stomp1.Stomp5 != stomp2.Stomp5 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp5 + "; was " + stomp2.Stomp5)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3, stomp2.Slot3) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot4, stomp2.Slot4) {
					slot1, _ := xml.Marshal(stomp1.Slot4)
					slot2, _ := xml.Marshal(stomp2.Slot4)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot5, stomp2.Slot5) {
					slot1, _ := xml.Marshal(stomp1.Slot5)
					slot2, _ := xml.Marshal(stomp2.Slot5)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite StompA2",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"StompA2",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA2
				stomp2 := preset2.StompA2
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if stomp1.Stomp4 != stomp2.Stomp4 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp4 + "; was " + stomp2.Stomp4)
				}
				if stomp1.Stomp5 != stomp2.Stomp5 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp5 + "; was " + stomp2.Stomp5)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3, stomp2.Slot3) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot4, stomp2.Slot4) {
					slot1, _ := xml.Marshal(stomp1.Slot4)
					slot2, _ := xml.Marshal(stomp2.Slot4)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot5, stomp2.Slot5) {
					slot1, _ := xml.Marshal(stomp1.Slot5)
					slot2, _ := xml.Marshal(stomp2.Slot5)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite StompB1",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"StompB1",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompB1
				stomp2 := preset2.StompB1
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if stomp1.Stomp4 != stomp2.Stomp4 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp4 + "; was " + stomp2.Stomp4)
				}
				if stomp1.Stomp5 != stomp2.Stomp5 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp5 + "; was " + stomp2.Stomp5)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3, stomp2.Slot3) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot4, stomp2.Slot4) {
					slot1, _ := xml.Marshal(stomp1.Slot4)
					slot2, _ := xml.Marshal(stomp2.Slot4)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot5, stomp2.Slot5) {
					slot1, _ := xml.Marshal(stomp1.Slot5)
					slot2, _ := xml.Marshal(stomp2.Slot5)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite StompB2",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"StompB2",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompB2
				stomp2 := preset2.StompB2
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if stomp1.Stomp4 != stomp2.Stomp4 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp4 + "; was " + stomp2.Stomp4)
				}
				if stomp1.Stomp5 != stomp2.Stomp5 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp5 + "; was " + stomp2.Stomp5)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3, stomp2.Slot3) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot4, stomp2.Slot4) {
					slot1, _ := xml.Marshal(stomp1.Slot4)
					slot2, _ := xml.Marshal(stomp2.Slot4)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot5, stomp2.Slot5) {
					slot1, _ := xml.Marshal(stomp1.Slot5)
					slot2, _ := xml.Marshal(stomp2.Slot5)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite StompB3",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"StompB3",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompB3
				stomp2 := preset2.StompB3
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if stomp1.Stomp4 != stomp2.Stomp4 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp4 + "; was " + stomp2.Stomp4)
				}
				if stomp1.Stomp5 != stomp2.Stomp5 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp5 + "; was " + stomp2.Stomp5)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3, stomp2.Slot3) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot4, stomp2.Slot4) {
					slot1, _ := xml.Marshal(stomp1.Slot4)
					slot2, _ := xml.Marshal(stomp2.Slot4)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot5, stomp2.Slot5) {
					slot1, _ := xml.Marshal(stomp1.Slot5)
					slot2, _ := xml.Marshal(stomp2.Slot5)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite StompStereo",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"StompStereo",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompStereo
				stomp2 := preset2.StompStereo
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite LoopFxA",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"LoopFxA",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.LoopFxA
				stomp2 := preset2.LoopFxA
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3, stomp2.Slot3) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite LoopFxB",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"LoopFxB",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.LoopFxB
				stomp2 := preset2.LoopFxB
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3, stomp2.Slot3) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite LoopFxC",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"LoopFxC",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.LoopFxC
				stomp2 := preset2.LoopFxC
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3, stomp2.Slot3) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite RackA",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"RackA",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.RackA
				stomp2 := preset2.RackA
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite RackB",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"RackB",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.RackB
				stomp2 := preset2.RackB
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite RackC",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"RackC",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.RackC
				stomp2 := preset2.RackC
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite RackDI",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"RackDI",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.RackDI
				stomp2 := preset2.RackDI
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite RackMaster",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"RackMaster",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.RackMaster
				stomp2 := preset2.RackMaster
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if stomp1.Stomp4 != stomp2.Stomp4 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp4 + "; was " + stomp2.Stomp4)
				}
				if stomp1.Stomp5 != stomp2.Stomp5 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp5 + "; was " + stomp2.Stomp5)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3, stomp2.Slot3) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot4, stomp2.Slot4) {
					slot1, _ := xml.Marshal(stomp1.Slot4)
					slot2, _ := xml.Marshal(stomp2.Slot4)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot5, stomp2.Slot5) {
					slot1, _ := xml.Marshal(stomp1.Slot5)
					slot2, _ := xml.Marshal(stomp2.Slot5)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite Stomps between types",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"StompA1",
				"StompA2",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA1
				stomp2 := preset2.StompA2
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if stomp1.Stomp4 != stomp2.Stomp4 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp4 + "; was " + stomp2.Stomp4)
				}
				if stomp1.Stomp5 != stomp2.Stomp5 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp5 + "; was " + stomp2.Stomp5)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3, stomp2.Slot3) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot4, stomp2.Slot4) {
					slot1, _ := xml.Marshal(stomp1.Slot4)
					slot2, _ := xml.Marshal(stomp2.Slot4)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot5, stomp2.Slot5) {
					slot1, _ := xml.Marshal(stomp1.Slot5)
					slot2, _ := xml.Marshal(stomp2.Slot5)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite stomp to stomp with fewer slots",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"StompA1",
				"StompStereo",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA1
				stomp2 := preset2.StompStereo
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Overwrite stomps to stomp with more slots",
			Command: "cpg",
			Args: []string{
				"-o",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"StompStereo",
				"StompA1",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompStereo
				stomp2 := preset2.StompA1
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp2.Stomp3 != EmptySlotGUID {
					return errors.New("unexpected GUID; expected " + EmptySlotGUID + "; was " + stomp2.Stomp3)
				}
				if stomp2.Stomp4 != EmptySlotGUID {
					return errors.New("unexpected GUID; expected " + EmptySlotGUID + "; was " + stomp2.Stomp4)
				}
				if stomp2.Stomp5 != EmptySlotGUID {
					return errors.New("unexpected GUID; expected " + EmptySlotGUID + "; was " + stomp2.Stomp5)
				}
				if !reflect.DeepEqual(stomp1.Slot0, stomp2.Slot0) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1, stomp2.Slot1) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2, stomp2.Slot2) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(emptySlot3(), stomp2.Slot3) {
					slot1, _ := xml.Marshal(Slot3{})
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(emptySlot4(), stomp2.Slot4) {
					slot1, _ := xml.Marshal(Slot4{})
					slot2, _ := xml.Marshal(stomp2.Slot4)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(emptySlot5(), stomp2.Slot5) {
					slot1, _ := xml.Marshal(Slot5{})
					slot2, _ := xml.Marshal(stomp2.Slot5)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Copy one amp to incompatible slot",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				"AmpA",
				"StompA1",
			},
			ExpectedError: "incompatible slots",
		},
		{
			Name:    "Don't copy amp when target is folder without recursive options",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
				"AmpA",
			},
			Expected: "-r not specified; omitting directory",
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension)
				file3 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension)
				file4 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				data3, err := ioutil.ReadFile(file3)
				data4, err := ioutil.ReadFile(file4)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				var preset3 PresetXMLV5
				var preset4 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				err = xml.Unmarshal(data3, &preset3)
				err = xml.Unmarshal(data4, &preset4)
				if err != nil {
					return err
				}
				if preset1.AmpA.Model == preset2.AmpA.Model {
					return errors.New("amps match; expected them to not")
				}
				if preset1.AmpA.Model == preset3.AmpA.Model {
					return errors.New("amps match; expected them to not")
				}
				if preset1.AmpA.Model == preset4.AmpA.Model {
					return errors.New("amps match; expected them to not")
				}
				return nil
			},
		},
		{
			Name:    "Copy gear to multiple presets recursive",
			Command: "cpg",
			Args: []string{
				"-r",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube"),
				"AmpA",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension)
				file3 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension)
				file4 := filepath.Join(workingDir, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				data3, err := ioutil.ReadFile(file3)
				data4, err := ioutil.ReadFile(file4)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				var preset3 PresetXMLV5
				var preset4 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				err = xml.Unmarshal(data3, &preset3)
				err = xml.Unmarshal(data4, &preset4)
				if err != nil {
					return err
				}
				if preset1.AmpA.Model != preset2.AmpA.Model {
					return errors.New("amps do not match; expected " + preset1.AmpA.Model + "; was " + preset2.AmpA.Model)
				}
				if preset1.AmpA.Model != preset3.AmpA.Model {
					return errors.New("amps do not match; expected " + preset1.AmpA.Model + "; was " + preset3.AmpA.Model)
				}
				if preset1.AmpA.Model != preset4.AmpA.Model {
					return errors.New("amps do not match; expected " + preset1.AmpA.Model + "; was " + preset4.AmpA.Model)
				}
				return nil
			},
		},
		{
			Name:    "Copy all amps",
			Command: "cpg",
			Args: []string{
				"-aa",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.AmpA.Model != preset2.AmpA.Model {
					return errors.New("amps do not match; expected " + preset1.AmpA.Model + "; was " + preset2.AmpA.Model)
				}
				if preset1.AmpB.Model != preset2.AmpB.Model {
					return errors.New("amps do not match; expected " + preset1.AmpB.Model + "; was " + preset2.AmpB.Model)
				}
				if preset1.AmpC.Model != preset2.AmpC.Model {
					return errors.New("amps do not match; expected " + preset1.AmpC.Model + "; was " + preset2.AmpC.Model)
				}
				return nil
			},
		},
		{
			Name:    "Copy all cabs",
			Command: "cpg",
			Args: []string{
				"-ac",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.CabA.CabModel != preset2.CabA.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabA.CabModel + "; was " + preset2.CabA.CabModel)
				}
				if preset1.CabB.CabModel != preset2.CabB.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabB.CabModel + "; was " + preset2.CabB.CabModel)
				}
				if preset1.CabC.CabModel != preset2.CabC.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabC.CabModel + "; was " + preset2.CabC.CabModel)
				}
				return nil
			},
		},
		{
			Name:    "Copy all amps and cabs",
			Command: "cpg",
			Args: []string{
				"-aa",
				"-ac",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.AmpA.Model != preset2.AmpA.Model {
					return errors.New("amps do not match; expected " + preset1.AmpA.Model + "; was " + preset2.AmpA.Model)
				}
				if preset1.AmpB.Model != preset2.AmpB.Model {
					return errors.New("amps do not match; expected " + preset1.AmpB.Model + "; was " + preset2.AmpB.Model)
				}
				if preset1.AmpC.Model != preset2.AmpC.Model {
					return errors.New("amps do not match; expected " + preset1.AmpC.Model + "; was " + preset2.AmpC.Model)
				}
				if preset1.CabA.CabModel != preset2.CabA.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabA.CabModel + "; was " + preset2.CabA.CabModel)
				}
				if preset1.CabB.CabModel != preset2.CabB.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabB.CabModel + "; was " + preset2.CabB.CabModel)
				}
				if preset1.CabC.CabModel != preset2.CabC.CabModel {
					return errors.New("cabs do not match; expected " + preset1.CabC.CabModel + "; was " + preset2.CabC.CabModel)
				}
				return nil
			},
		},
		{
			Name:    "Copy all fx",
			Command: "cpg",
			Args: []string{
				"-ae",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "Default"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				if preset1.StompA1.Stomp0 != preset2.StompA1.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.StompA1.Stomp0 + "; was " + preset2.StompA1.Stomp0)
				}
				if preset1.StompA1.Stomp1 != preset2.StompA1.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.StompA1.Stomp1 + "; was " + preset2.StompA1.Stomp1)
				}
				if preset1.StompA1.Stomp2 != preset2.StompA1.Stomp2 {
					return errors.New("stomps do not match; expected " + preset1.StompA1.Stomp2 + "; was " + preset2.StompA1.Stomp2)
				}
				if preset1.StompA1.Stomp3 != preset2.StompA1.Stomp3 {
					return errors.New("stomps do not match; expected " + preset1.StompA1.Stomp3 + "; was " + preset2.StompA1.Stomp3)
				}
				if preset1.StompA1.Stomp4 != preset2.StompA1.Stomp4 {
					return errors.New("stomps do not match; expected " + preset1.StompA1.Stomp4 + "; was " + preset2.StompA1.Stomp4)
				}
				if preset1.StompA1.Stomp5 != preset2.StompA1.Stomp5 {
					return errors.New("stomps do not match; expected " + preset1.StompA1.Stomp5 + "; was " + preset2.StompA1.Stomp5)
				}
				if preset1.StompA2.Stomp0 != preset2.StompA2.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.StompA2.Stomp0 + "; was " + preset2.StompA2.Stomp0)
				}
				if preset1.StompA2.Stomp1 != preset2.StompA2.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.StompA2.Stomp1 + "; was " + preset2.StompA2.Stomp1)
				}
				if preset1.StompA2.Stomp2 != preset2.StompA2.Stomp2 {
					return errors.New("stomps do not match; expected " + preset1.StompA2.Stomp2 + "; was " + preset2.StompA2.Stomp2)
				}
				if preset1.StompA2.Stomp3 != preset2.StompA2.Stomp3 {
					return errors.New("stomps do not match; expected " + preset1.StompA2.Stomp3 + "; was " + preset2.StompA2.Stomp3)
				}
				if preset1.StompA2.Stomp4 != preset2.StompA2.Stomp4 {
					return errors.New("stomps do not match; expected " + preset1.StompA2.Stomp4 + "; was " + preset2.StompA2.Stomp4)
				}
				if preset1.StompA2.Stomp5 != preset2.StompA2.Stomp5 {
					return errors.New("stomps do not match; expected " + preset1.StompA2.Stomp5 + "; was " + preset2.StompA2.Stomp5)
				}
				if preset1.StompStereo.Stomp0 != preset2.StompStereo.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.StompStereo.Stomp0 + "; was " + preset2.StompStereo.Stomp0)
				}
				if preset1.StompStereo.Stomp1 != preset2.StompStereo.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.StompStereo.Stomp1 + "; was " + preset2.StompStereo.Stomp1)
				}
				if preset1.StompStereo.Stomp2 != preset2.StompStereo.Stomp2 {
					return errors.New("stomps do not match; expected " + preset1.StompStereo.Stomp2 + "; was " + preset2.StompStereo.Stomp2)
				}
				if preset1.StompB1.Stomp0 != preset2.StompB1.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.StompB1.Stomp0 + "; was " + preset2.StompB1.Stomp0)
				}
				if preset1.StompB1.Stomp1 != preset2.StompB1.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.StompB1.Stomp1 + "; was " + preset2.StompB1.Stomp1)
				}
				if preset1.StompB1.Stomp2 != preset2.StompB1.Stomp2 {
					return errors.New("stomps do not match; expected " + preset1.StompB1.Stomp2 + "; was " + preset2.StompB1.Stomp2)
				}
				if preset1.StompB1.Stomp3 != preset2.StompB1.Stomp3 {
					return errors.New("stomps do not match; expected " + preset1.StompB1.Stomp3 + "; was " + preset2.StompB1.Stomp3)
				}
				if preset1.StompB1.Stomp4 != preset2.StompB1.Stomp4 {
					return errors.New("stomps do not match; expected " + preset1.StompB1.Stomp4 + "; was " + preset2.StompB1.Stomp4)
				}
				if preset1.StompB1.Stomp5 != preset2.StompB1.Stomp5 {
					return errors.New("stomps do not match; expected " + preset1.StompB1.Stomp5 + "; was " + preset2.StompB1.Stomp5)
				}
				if preset1.StompB2.Stomp0 != preset2.StompB2.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.StompB2.Stomp0 + "; was " + preset2.StompB2.Stomp0)
				}
				if preset1.StompB2.Stomp1 != preset2.StompB2.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.StompB2.Stomp1 + "; was " + preset2.StompB2.Stomp1)
				}
				if preset1.StompB2.Stomp2 != preset2.StompB2.Stomp2 {
					return errors.New("stomps do not match; expected " + preset1.StompB2.Stomp2 + "; was " + preset2.StompB2.Stomp2)
				}
				if preset1.StompB2.Stomp3 != preset2.StompB2.Stomp3 {
					return errors.New("stomps do not match; expected " + preset1.StompB2.Stomp3 + "; was " + preset2.StompB2.Stomp3)
				}
				if preset1.StompB2.Stomp4 != preset2.StompB2.Stomp4 {
					return errors.New("stomps do not match; expected " + preset1.StompB2.Stomp4 + "; was " + preset2.StompB2.Stomp4)
				}
				if preset1.StompB2.Stomp5 != preset2.StompB2.Stomp5 {
					return errors.New("stomps do not match; expected " + preset1.StompB2.Stomp5 + "; was " + preset2.StompB2.Stomp5)
				}
				if preset1.StompB3.Stomp0 != preset2.StompB3.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.StompB3.Stomp0 + "; was " + preset2.StompB3.Stomp0)
				}
				if preset1.StompB3.Stomp1 != preset2.StompB3.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.StompB3.Stomp1 + "; was " + preset2.StompB3.Stomp1)
				}
				if preset1.StompB3.Stomp2 != preset2.StompB3.Stomp2 {
					return errors.New("stomps do not match; expected " + preset1.StompB3.Stomp2 + "; was " + preset2.StompB3.Stomp2)
				}
				if preset1.StompB3.Stomp3 != preset2.StompB3.Stomp3 {
					return errors.New("stomps do not match; expected " + preset1.StompB3.Stomp3 + "; was " + preset2.StompB3.Stomp3)
				}
				if preset1.StompB3.Stomp4 != preset2.StompB3.Stomp4 {
					return errors.New("stomps do not match; expected " + preset1.StompB3.Stomp4 + "; was " + preset2.StompB3.Stomp4)
				}
				if preset1.StompB3.Stomp5 != preset2.StompB3.Stomp5 {
					return errors.New("stomps do not match; expected " + preset1.StompB3.Stomp5 + "; was " + preset2.StompB3.Stomp5)
				}
				if preset1.LoopFxA.Stomp0 != preset2.LoopFxA.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxA.Stomp0 + "; was " + preset2.LoopFxA.Stomp0)
				}
				if preset1.LoopFxA.Stomp1 != preset2.LoopFxA.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxA.Stomp1 + "; was " + preset2.LoopFxA.Stomp1)
				}
				if preset1.LoopFxA.Stomp2 != preset2.LoopFxA.Stomp2 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxA.Stomp2 + "; was " + preset2.LoopFxA.Stomp2)
				}
				if preset1.LoopFxA.Stomp3 != preset2.LoopFxA.Stomp3 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxA.Stomp3 + "; was " + preset2.LoopFxA.Stomp3)
				}
				if preset1.LoopFxB.Stomp0 != preset2.LoopFxB.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxB.Stomp0 + "; was " + preset2.LoopFxB.Stomp0)
				}
				if preset1.LoopFxB.Stomp1 != preset2.LoopFxB.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxB.Stomp1 + "; was " + preset2.LoopFxB.Stomp1)
				}
				if preset1.LoopFxB.Stomp2 != preset2.LoopFxB.Stomp2 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxB.Stomp2 + "; was " + preset2.LoopFxB.Stomp2)
				}
				if preset1.LoopFxB.Stomp3 != preset2.LoopFxB.Stomp3 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxB.Stomp3 + "; was " + preset2.LoopFxB.Stomp3)
				}
				if preset1.LoopFxC.Stomp0 != preset2.LoopFxC.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxC.Stomp0 + "; was " + preset2.LoopFxC.Stomp0)
				}
				if preset1.LoopFxC.Stomp1 != preset2.LoopFxC.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxC.Stomp1 + "; was " + preset2.LoopFxC.Stomp1)
				}
				if preset1.LoopFxC.Stomp2 != preset2.LoopFxC.Stomp2 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxC.Stomp2 + "; was " + preset2.LoopFxC.Stomp2)
				}
				if preset1.LoopFxC.Stomp3 != preset2.LoopFxC.Stomp3 {
					return errors.New("stomps do not match; expected " + preset1.LoopFxC.Stomp3 + "; was " + preset2.LoopFxC.Stomp3)
				}
				if preset1.RackA.Stomp0 != preset2.RackA.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.RackA.Stomp0 + "; was " + preset2.RackA.Stomp0)
				}
				if preset1.RackA.Stomp1 != preset2.RackA.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.RackA.Stomp1 + "; was " + preset2.RackA.Stomp1)
				}
				if preset1.RackB.Stomp0 != preset2.RackB.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.RackB.Stomp0 + "; was " + preset2.RackB.Stomp0)
				}
				if preset1.RackB.Stomp1 != preset2.RackB.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.RackB.Stomp1 + "; was " + preset2.RackB.Stomp1)
				}
				if preset1.RackC.Stomp0 != preset2.RackC.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.RackC.Stomp0 + "; was " + preset2.RackC.Stomp0)
				}
				if preset1.RackC.Stomp1 != preset2.RackC.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.RackC.Stomp1 + "; was " + preset2.RackC.Stomp1)
				}
				if preset1.RackMaster.Stomp0 != preset2.RackMaster.Stomp0 {
					return errors.New("stomps do not match; expected " + preset1.RackMaster.Stomp0 + "; was " + preset2.RackMaster.Stomp0)
				}
				if preset1.RackMaster.Stomp1 != preset2.RackMaster.Stomp1 {
					return errors.New("stomps do not match; expected " + preset1.RackMaster.Stomp1 + "; was " + preset2.RackMaster.Stomp1)
				}
				if preset1.RackMaster.Stomp2 != preset2.RackMaster.Stomp2 {
					return errors.New("stomps do not match; expected " + preset1.RackMaster.Stomp2 + "; was " + preset2.RackMaster.Stomp2)
				}
				if preset1.RackMaster.Stomp3 != preset2.RackMaster.Stomp3 {
					return errors.New("stomps do not match; expected " + preset1.RackMaster.Stomp3 + "; was " + preset2.RackMaster.Stomp3)
				}
				if preset1.RackMaster.Stomp4 != preset2.RackMaster.Stomp4 {
					return errors.New("stomps do not match; expected " + preset1.RackMaster.Stomp4 + "; was " + preset2.RackMaster.Stomp4)
				}
				if preset1.RackMaster.Stomp5 != preset2.RackMaster.Stomp5 {
					return errors.New("stomps do not match; expected " + preset1.RackMaster.Stomp5 + "; was " + preset2.RackMaster.Stomp5)
				}
				return nil
			},
		},
		{
			Name:    "Append fx",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearEmpty"+PresetExtension),
				"StompA1",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearEmpty"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA1
				stomp2 := preset2.StompA1
				if stomp1.Stomp0 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp1)
				}
				if !reflect.DeepEqual(stomp1.Slot0.Attrs, stomp2.Slot1.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Append fx different types",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearEmpty"+PresetExtension),
				"StompA2",
				"StompA1",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearEmpty"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA2
				stomp2 := preset2.StompA1
				if stomp1.Stomp0 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp1)
				}
				if !reflect.DeepEqual(stomp1.Slot0.Attrs, stomp2.Slot1.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Append fx with truncation",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearEmpty"+PresetExtension),
				"StompB1",
				"StompStereo",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearEmpty"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompB1
				stomp2 := preset2.StompStereo
				if stomp1.Stomp0 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp1)
				}
				if stomp1.Stomp1 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp1 + "; was " + stomp2.Stomp2)
				}
				if !reflect.DeepEqual(stomp1.Slot0.Attrs, stomp2.Slot1.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1.Attrs, stomp2.Slot2.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Insert fx",
			Command: "cpg",
			Args: []string{
				"-i",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearEmpty"+PresetExtension),
				"StompA2",
				"StompA1",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearEmpty"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA2
				stomp2 := preset2.StompA1
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp2.Stomp1 != "a1000000-0000-0000-0000-000000000000" {
					return errors.New("stomps do not match; expected a1000000-0000-0000-0000-000000000000; was " + stomp2.Stomp1)
				}
				if !reflect.DeepEqual(stomp1.Slot0.Attrs, stomp2.Slot0.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Check order of existing pedals after append",
			Command: "cpg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSparseSource2"+PresetExtension),
				"StompA2",
				"StompB1",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSparseSource2"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA2
				stomp1b := preset1.StompB1
				stomp2 := preset2.StompB1
				if stomp1.Stomp0 != stomp2.Stomp4 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp4)
				}
				if stomp1b.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1b.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1b.Stomp1 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1b.Stomp1 + "; was " + stomp2.Stomp1)
				}
				if stomp1b.Stomp2 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1b.Stomp2 + "; was " + stomp2.Stomp2)
				}
				if stomp1b.Stomp3 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1b.Stomp3 + "; was " + stomp2.Stomp3)
				}
				if !reflect.DeepEqual(stomp1.Slot0.Attrs, stomp2.Slot4.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot4)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1b.Slot0.Attrs, stomp2.Slot0.Attrs) {
					slot1, _ := xml.Marshal(stomp1b.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1b.Slot1.Attrs, stomp2.Slot1.Attrs) {
					slot1, _ := xml.Marshal(stomp1b.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1b.Slot2.Attrs, stomp2.Slot2.Attrs) {
					slot1, _ := xml.Marshal(stomp1b.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1b.Slot3.Attrs, stomp2.Slot3.Attrs) {
					slot1, _ := xml.Marshal(stomp1b.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Check order of existing pedals after insert",
			Command: "cpg",
			Args: []string{
				"-i",
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSparseSource2"+PresetExtension),
				"StompA2",
				"StompB1",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension)
				file2 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSparseSource2"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				data2, err := ioutil.ReadFile(file2)
				var preset1 PresetXMLV5
				var preset2 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				err = xml.Unmarshal(data2, &preset2)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA2
				stomp1b := preset1.StompB1
				stomp2 := preset2.StompB1
				if stomp1.Stomp0 != stomp2.Stomp0 {
					return errors.New("stomps do not match; expected " + stomp1.Stomp0 + "; was " + stomp2.Stomp0)
				}
				if stomp1b.Stomp0 != stomp2.Stomp1 {
					return errors.New("stomps do not match; expected " + stomp1b.Stomp0 + "; was " + stomp2.Stomp1)
				}
				if stomp1b.Stomp1 != stomp2.Stomp2 {
					return errors.New("stomps do not match; expected " + stomp1b.Stomp1 + "; was " + stomp2.Stomp2)
				}
				if stomp1b.Stomp2 != stomp2.Stomp3 {
					return errors.New("stomps do not match; expected " + stomp1b.Stomp2 + "; was " + stomp2.Stomp3)
				}
				if stomp1b.Stomp3 != stomp2.Stomp4 {
					return errors.New("stomps do not match; expected " + stomp1b.Stomp3 + "; was " + stomp2.Stomp4)
				}
				if !reflect.DeepEqual(stomp1.Slot0.Attrs, stomp2.Slot0.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot0)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1b.Slot0.Attrs, stomp2.Slot1.Attrs) {
					slot1, _ := xml.Marshal(stomp1b.Slot0)
					slot2, _ := xml.Marshal(stomp2.Slot1)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1b.Slot1.Attrs, stomp2.Slot2.Attrs) {
					slot1, _ := xml.Marshal(stomp1b.Slot1)
					slot2, _ := xml.Marshal(stomp2.Slot2)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1b.Slot2.Attrs, stomp2.Slot3.Attrs) {
					slot1, _ := xml.Marshal(stomp1b.Slot2)
					slot2, _ := xml.Marshal(stomp2.Slot3)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1b.Slot3.Attrs, stomp2.Slot4.Attrs) {
					slot1, _ := xml.Marshal(stomp1b.Slot3)
					slot2, _ := xml.Marshal(stomp2.Slot4)
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		// TODO: copy gear between v4 an v5 presets
		{
			Name:    "Remove fx",
			Command: "rmg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension),
				"StompA1",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				var preset1 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA1
				if stomp1.Stomp0 != EmptySlotGUID {
					return errors.New("stomps do not match; expected " + EmptySlotGUID + "; was " + stomp1.Stomp0)
				}
				if stomp1.Stomp1 != EmptySlotGUID {
					return errors.New("stomps do not match; expected " + EmptySlotGUID + "; was " + stomp1.Stomp1)
				}
				if stomp1.Stomp2 != EmptySlotGUID {
					return errors.New("stomps do not match; expected " + EmptySlotGUID + "; was " + stomp1.Stomp2)
				}
				if stomp1.Stomp3 != EmptySlotGUID {
					return errors.New("stomps do not match; expected " + EmptySlotGUID + "; was " + stomp1.Stomp3)
				}
				if stomp1.Stomp4 != EmptySlotGUID {
					return errors.New("stomps do not match; expected " + EmptySlotGUID + "; was " + stomp1.Stomp4)
				}
				if stomp1.Stomp5 != EmptySlotGUID {
					return errors.New("stomps do not match; expected " + EmptySlotGUID + "; was " + stomp1.Stomp5)
				}
				if !reflect.DeepEqual(stomp1.Slot0.Attrs, Slot0{}.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(Slot0{})
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot1.Attrs, Slot1{}.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot1)
					slot2, _ := xml.Marshal(Slot1{})
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot2.Attrs, Slot2{}.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot2)
					slot2, _ := xml.Marshal(Slot2{})
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot3.Attrs, Slot3{}.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot3)
					slot2, _ := xml.Marshal(Slot3{})
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot4.Attrs, Slot4{}.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot4)
					slot2, _ := xml.Marshal(Slot4{})
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if !reflect.DeepEqual(stomp1.Slot5.Attrs, Slot5{}.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot5)
					slot2, _ := xml.Marshal(Slot5{})
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				return nil
			},
		},
		{
			Name:    "Remove fx one slot",
			Command: "rmg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSource"+PresetExtension),
				"StompA1",
				"Slot0",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSource"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				var preset1 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA1
				if stomp1.Stomp0 != EmptySlotGUID {
					return errors.New("stomps do not match; expected " + EmptySlotGUID + "; was " + stomp1.Stomp0)
				}
				if stomp1.Stomp1 == EmptySlotGUID {
					return errors.New("stomps do not match; expected to not be empty")
				}
				if stomp1.Stomp2 == EmptySlotGUID {
					return errors.New("stomps do not match; expected to not be empty")
				}
				if stomp1.Stomp3 == EmptySlotGUID {
					return errors.New("stomps do not match; expected to not be empty")
				}
				if stomp1.Stomp4 == EmptySlotGUID {
					return errors.New("stomps do not match; expected to not be empty")
				}
				if stomp1.Stomp5 == EmptySlotGUID {
					return errors.New("stomps do not match; expected to not be empty")
				}
				if !reflect.DeepEqual(stomp1.Slot0.Attrs, Slot0{}.Attrs) {
					slot1, _ := xml.Marshal(stomp1.Slot0)
					slot2, _ := xml.Marshal(Slot0{})
					return errors.New("stomps do not match; expected " + string(selfClose(slot1)) + "; was " + string(selfClose(slot2)))
				}
				if reflect.DeepEqual(stomp1.Slot1.Attrs, Slot1{}.Attrs) {
					return errors.New("stomps do not match; expected to not be empty")
				}
				if reflect.DeepEqual(stomp1.Slot2.Attrs, Slot2{}.Attrs) {
					return errors.New("stomps do not match; expected to not be empty")
				}
				if reflect.DeepEqual(stomp1.Slot3.Attrs, Slot3{}.Attrs) {
					return errors.New("stomps do not match; expected to not be empty")
				}
				if reflect.DeepEqual(stomp1.Slot4.Attrs, Slot4{}.Attrs) {
					return errors.New("stomps do not match; expected to not be empty")
				}
				if reflect.DeepEqual(stomp1.Slot5.Attrs, Slot5{}.Attrs) {
					return errors.New("stomps do not match; expected to not be empty")
				}
				return nil
			},
		},
		{
			Name:    "Set slot parameter",
			Command: "sg",
			Args: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension),
				"Preset.StompA1.Slot0.Setting_1=1",
			},
			CustomAssertion: func(workingDir string) error {
				file1 := filepath.Join(workingDir, PresetsFolder, "Amps", "TestGearSparseSource"+PresetExtension)
				data1, err := ioutil.ReadFile(file1)
				var preset1 PresetXMLV5
				err = xml.Unmarshal(data1, &preset1)
				if err != nil {
					return err
				}
				stomp1 := preset1.StompA1
				for _, attr := range stomp1.Slot0.Attrs {
					if attr.Name.Local == "Setting_1" && attr.Value != "1" {
						return errors.New("slot setting not updated; expected 1; actual " + attr.Value)
					}
				}
				return nil
			},
		},
		{
			Name:         "Import first arg is a profile directory",
			Command:      "import",
			WorkDirCount: 2,
			Args: []string{
				".",
				TestDataRoot,
			},
			ExpectedError: "both args must be the root of an Amplitube profile",
		},
		{
			Name:         "Import second arg is a profile directory",
			Command:      "import",
			WorkDirCount: 2,
			Args: []string{
				TestDataRoot,
				".",
			},
			ExpectedError: "both args must be the root of an Amplitube profile",
		},
		{
			Name:         "Import args must be profile directories",
			Command:      "import",
			WorkDirCount: 2,
			Args: []string{
				".",
				".",
			},
			ExpectedError: "both args must be the root of an Amplitube profile",
		},
		{
			Name:    "Import args cannot be the same profile",
			Command: "import",
			Args: []string{
				TestDataRoot,
				TestDataRoot,
			},
			ExpectedError: "args cannot reference the same Amplitube profile",
		},
		{
			Name:         "Import presets from separate profile",
			Command:      "import",
			WorkDirCount: 2,
			Args: []string{
				TestDataRoot + "[1]",
				TestDataRoot,
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "THD", "BiValve"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Empty"),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "THD"),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "THD", "BiValve"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:         "Import arg referencing preset resolves to containing profile",
			Command:      "import",
			WorkDirCount: 2,
			Args: []string{
				filepath.Join(TestDataRoot+"[1]", PresetsFolder, "Amps"),
				TestDataRoot,
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "THD", "BiValve"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Empty"),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "THD"),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "THD", "BiValve"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Amps2", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		{
			Name:    "Import from Amplitube 4 Profile",
			Command: "import",
			Args: []string{
				"amp4data",
				TestDataRoot,
			},
			ExpectExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Default"+PresetExtension4),
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Import", "Default"+PresetExtension4),
			},
		},
		{
			Name:          "Reindex arg required",
			Command:       "reindex",
			Args:          []string{},
			ExpectedError: "arg must be the root of an Amplitube profile",
		},
		{
			Name:    "Reindex required profile",
			Command: "reindex",
			Args: []string{
				".",
			},
			ExpectedError: "arg must be the root of an Amplitube profile",
		},
		{
			Name:         "Reindex db file paths",
			Command:      "reindex",
			WorkDirCount: 2,
			Args: []string{
				TestDataRoot,
			},
			CustomSetup: func(workingDirs []string) {
				os.Rename(filepath.Join(workingDirs[1], "Presets.db"), filepath.Join(workingDirs[0], "Presets.db"))
			},
			ExpectDBExists: []string{
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps2", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps2", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot, PresetsFolder, "Amps2", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
			ExpectDBNotExist: []string{
				filepath.Join(TestDataRoot+"[1]", PresetsFolder, "Amps", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot+"[1]", PresetsFolder, "Amps", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot+"[1]", PresetsFolder, "Amps", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
				filepath.Join(TestDataRoot+"[1]", PresetsFolder, "Amps", "THD", "BiValve"+PresetExtension),
				filepath.Join(TestDataRoot+"[1]", PresetsFolder, "Amps", "Default"+PresetExtension),
				filepath.Join(TestDataRoot+"[1]", PresetsFolder, "Amps2", "Amplitube", "American Tube Clean 1"+PresetExtension),
				filepath.Join(TestDataRoot+"[1]", PresetsFolder, "Amps2", "Amplitube", "Metal", "Metal Clean T"+PresetExtension),
				filepath.Join(TestDataRoot+"[1]", PresetsFolder, "Amps2", "Amplitube", "SVX", "SVX-4B"+PresetExtension),
			},
		},
		// TODO: remove orphans and add missing db records on reindex
	} {
		t.Run(tc.Name, func(t *testing.T) {

			workDirCount := tc.WorkDirCount

			if workDirCount == 0 {
				workDirCount = 1
			}

			var workingDirs []string

			for i := 0; i < workDirCount; i++ {
				workingDir := setupData()
				defer cleanUpData(workingDir)
				workingDirs = append(workingDirs, workingDir)
			}

			if tc.CustomSetup != nil {
				tc.CustomSetup(workingDirs)
			}

			var filteredArgs []string

			for _, v := range tc.Args {
				filteredArgs = append(filteredArgs, rehomePath(v, workingDirs))
			}

			out = bytes.NewBuffer(nil)
			err := ExecuteCommand(tc.Command, filteredArgs)

			if tc.Verbose {
				t.Logf("%s", out.(*bytes.Buffer).String())
			}

			if err != nil && !strings.Contains(err.Error(), tc.ExpectedError) {
				t.Fatal(err)
			}

			if err == nil && tc.ExpectedError != "" {
				t.Errorf("Expected error '%s'; there was none", tc.ExpectedError)
			}

			if err != nil && tc.ExpectedError == "" {
				t.Errorf("Error was '%s'; no error was expected", err)
			}

			if !strings.Contains(out.(*bytes.Buffer).String(), tc.Expected) {
				t.Errorf("wanted '%s'; was '%s'", tc.Expected, out.(*bytes.Buffer).String())
			}

			for _, path := range tc.ExpectExists {
				path = rehomePath(path, workingDirs)
				if info, _ := os.Stat(path); info == nil {
					t.Errorf("Expected file '%s' to exist.", path)
				}
			}

			for _, path := range tc.ExpectNotExist {
				path = rehomePath(path, workingDirs)
				if info, _ := os.Stat(path); info != nil {
					t.Errorf("Expected file '%s' to not exists.", path)
				}
			}

			database, err := sql.Open("sqlite3", strings.ReplaceAll(filepath.Join(workingDirs[0], "Presets.db"), "\\", "/"))

			if err != nil {
				log.Fatal("Failed validating result: ", err)
			}

			statement, err := database.Prepare("select count(id) from pXcPresets where OriginalFileName = ?")

			if err != nil {
				log.Fatal("Failed validating result: ", err)
			}

			for _, path := range tc.ExpectDBExists {
				path = filepath.Join(filepath.Dir(workingDirs[0]), rehomePath(path, workingDirs))
				rs, _ := statement.Query(path)
				var count = -1
				for rs.Next() {
					rs.Scan(&count)
				}
				if count != 1 {
					t.Errorf("Incorrect number of database record of %d for file '%s'; expected 1", count, path)
				}
			}

			for _, path := range tc.ExpectDBNotExist {
				path = filepath.Join(filepath.Dir(workingDirs[0]), rehomePath(path, workingDirs))
				rs, _ := statement.Query(path)
				var count = -1
				for rs.Next() {
					rs.Scan(&count)
				}
				if count > 0 {
					t.Errorf("Incorrect number of database record of %d for file '%s'; expected 0", count, path)
				}
			}

			if tc.CustomAssertion != nil {
				if err = tc.CustomAssertion(workingDirs[0]); err != nil {
					t.Error(err)
				}
			}

			statement.Close()
			database.Close()

		})
	}

}

func setupData() string {

	toCopy := map[string]string{}

	randomUUID, _ := uuid.NewRandom()

	pwd, _ := os.Getwd()

	workingDir := filepath.Join(pwd, randomUUID.String())

	filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("Failed setting up test: ", err)
		}
		source := filepath.Join(pwd, path)
		if !info.IsDir() {
			toCopy[source] = strings.ReplaceAll(source, string(filepath.Separator)+"testdata"+string(filepath.Separator), string(filepath.Separator)+randomUUID.String()+string(filepath.Separator))
		}
		return nil
	})

	for source, target := range toCopy {
		os.MkdirAll(filepath.Dir(target), 0775)
		data, _ := ioutil.ReadFile(source)
		ioutil.WriteFile(target, data, 0664)
	}

	filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Base(path) == EmptyPlaceholder {
			os.Remove(path)
		}
		return nil
	})

	database, err := sql.Open("sqlite3", strings.ReplaceAll(filepath.Join(pwd, randomUUID.String(), "Presets.db"), "\\", "/"))

	if err != nil {
		log.Fatal("Failed setting up test: ", err)
	}

	tx, err := database.Begin()

	if err != nil {
		database.Close()
		log.Fatal("Failed setting up test: ", err)
	}

	statement, err := database.Prepare("update pXcPresets set OriginalFileName = ? || replace(OriginalFileName,'\\',?), FileFolder = ? || replace(FileFolder,'\\',?)")

	_, err = statement.Exec(filepath.Join(workingDir, PresetsFolder), string(filepath.Separator), filepath.Join(workingDir, PresetsFolder), string(filepath.Separator))

	if err != nil {
		tx.Rollback()
		database.Close()
		log.Fatal("Failed setting up test: ", err)
	}

	tx.Commit()
	statement.Close()
	database.Close()

	return workingDir
}

func cleanUpData(workingDir string) {
	os.RemoveAll(workingDir)
}

func rehomePath(path string, workingDirs []string) string {
	match, _ := regexp.MatchString(TestDataRoot+"\\[\\d+\\]", path)
	if match {
		index := path[len(TestDataRoot)+1 : strings.Index(path, "]")]
		i, _ := strconv.Atoi(index)
		return strings.ReplaceAll(path, TestDataRoot+"["+index+"]", filepath.Base(workingDirs[i]))
	} else {
		return strings.ReplaceAll(path, TestDataRoot, filepath.Base(workingDirs[0]))
	}
}

type TestSlotGroup3 struct {
	XMLName xml.Name  `xml:"Test"`
	Slot0   TestSlot0 `xml:""`
	Slot1   TestSlot1 `xml:""`
	Slot2   TestSlot2 `xml:""`
}

type TestSlotGroup1and5 struct {
	XMLName xml.Name   `xml:"Test"`
	Slot0   TestSlot0  `xml:""`
	Slot1   EmptySlot1 `xml:""`
	Slot2   EmptySlot2 `xml:""`
	Slot3   EmptySlot3 `xml:""`
	Slot4   EmptySlot4 `xml:""`
	Slot5   EmptySlot5 `xml:""`
}

type TestSlotGroup2and4 struct {
	XMLName xml.Name   `xml:"Test"`
	Slot0   TestSlot0  `xml:""`
	Slot1   TestSlot1  `xml:""`
	Slot2   EmptySlot2 `xml:""`
	Slot3   EmptySlot3 `xml:""`
	Slot4   EmptySlot4 `xml:""`
	Slot5   EmptySlot5 `xml:""`
}

type TestSlotGroup6 struct {
	XMLName xml.Name  `xml:"Test"`
	Slot0   TestSlot0 `xml:""`
	Slot1   TestSlot1 `xml:""`
	Slot2   TestSlot2 `xml:""`
	Slot3   TestSlot3 `xml:""`
	Slot4   TestSlot4 `xml:""`
	Slot5   TestSlot5 `xml:""`
}

type TestSlotGroup3And3 struct {
	XMLName xml.Name   `xml:"Test"`
	Slot0   TestSlot0  `xml:""`
	Slot1   TestSlot1  `xml:""`
	Slot2   TestSlot2  `xml:""`
	Slot3   EmptySlot3 `xml:""`
	Slot4   EmptySlot4 `xml:""`
	Slot5   EmptySlot5 `xml:""`
}

type TestSlot0 struct {
	XMLName xml.Name `xml:"Slot0"`
	Name    string   `xml:",attr"`
}

type TestSlot1 struct {
	XMLName xml.Name `xml:"Slot1"`
	Name    string   `xml:",attr"`
}

type TestSlot2 struct {
	XMLName xml.Name `xml:"Slot2"`
	Name    string   `xml:",attr"`
}

type TestSlot3 struct {
	XMLName xml.Name `xml:"Slot3"`
	Name    string   `xml:",attr"`
}

type TestSlot4 struct {
	XMLName xml.Name `xml:"Slot4"`
	Name    string   `xml:",attr"`
}

type TestSlot5 struct {
	XMLName xml.Name `xml:"Slot5"`
	Name    string   `xml:",attr"`
}

type EmptySlot0 struct {
	XMLName xml.Name `xml:"Slot0"`
}

type EmptySlot1 struct {
	XMLName xml.Name `xml:"Slot1"`
}

type EmptySlot2 struct {
	XMLName xml.Name `xml:"Slot2"`
}

type EmptySlot3 struct {
	XMLName xml.Name `xml:"Slot3"`
}

type EmptySlot4 struct {
	XMLName xml.Name `xml:"Slot4"`
}

type EmptySlot5 struct {
	XMLName xml.Name `xml:"Slot5"`
}

func emptySlot3() Slot3 {
	return Slot3{
		XMLName: xml.Name{
			Local: "Slot3",
		},
		Attrs: nil,
	}
}

func emptySlot4() Slot4 {
	return Slot4{
		XMLName: xml.Name{
			Local: "Slot4",
		},
		Attrs: nil,
	}
}

func emptySlot5() Slot5 {
	return Slot5{
		XMLName: xml.Name{
			Local: "Slot5",
		},
		Attrs: nil,
	}
}
