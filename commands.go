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
	"flag"
	"path/filepath"
)

type ExecutionContext struct {
	Profile  string
	Args     []string
	Options  map[string]interface{}
	Database *sql.DB
}

type Runner func(ExecutionContext) error
type DatabaseFactory func(ExecutionContext) (*sql.DB, error)

type Command struct {
	Flags           *flag.FlagSet
	Runner          Runner
	Options         map[string]interface{}
	DatabaseFactory DatabaseFactory
}

func defaultDatabaseFactory(context ExecutionContext) (*sql.DB, error) {
	if len(context.Args) == 0 {
		return nil, nil
	}
	source, err := filepath.Abs(context.Args[0])
	sourceProfile, err := resolveToProfile(source)
	if err != nil {
		return nil, nil
	}
	database, err := openDatabase(filepath.Join(sourceProfile, "Presets.db"))
	if err == nil {
		if !isV5Database(database) {
			return nil, errors.New("incompatible source database version")
		}
	}
	return database, err
}

func secondArgDatabaseFactory(context ExecutionContext) (*sql.DB, error) {
	if len(context.Args) == 0 {
		return nil, nil
	}
	source, err := filepath.Abs(context.Args[1])
	sourceProfile, err := resolveToProfile(source)
	if err != nil {
		return nil, nil
	}
	database, err := openDatabase(filepath.Join(sourceProfile, "Presets.db"))
	if err == nil {
		if !isV5Database(database) {
			return nil, errors.New("incompatible database version")
		}
	}
	return database, err
}

func nilDatabaseFactory(context ExecutionContext) (*sql.DB, error) {
	return nil, nil
}

func ExecuteCommand(cmd string, args []string) error {

	var lsFlags = flag.NewFlagSet("ls", flag.ExitOnError)
	var lsgFlags = flag.NewFlagSet("lsg", flag.ExitOnError)
	var mkDirFlags = flag.NewFlagSet("mkdir", flag.ExitOnError)
	var rmFlags = flag.NewFlagSet("rm", flag.ExitOnError)
	var rmgFlags = flag.NewFlagSet("rmg", flag.ExitOnError)
	var mvFlags = flag.NewFlagSet("mv", flag.ExitOnError)
	var cpFlags = flag.NewFlagSet("cp", flag.ExitOnError)
	var cpgFlags = flag.NewFlagSet("cpg", flag.ExitOnError)
	var importFlags = flag.NewFlagSet("import", flag.ExitOnError)
	var reindexFlags = flag.NewFlagSet("reindex", flag.ExitOnError)
	var sgFlags = flag.NewFlagSet("sg", flag.ExitOnError)

	var commands = map[string]*Command{
		"cp": {
			Flags:           cpFlags,
			Runner:          copy,
			DatabaseFactory: defaultDatabaseFactory,
			Options: map[string]interface{}{
				"recursive": cpFlags.Bool("r", false, "Copy to subfolders"),
			},
		},
		"cpg": {
			Flags:           cpgFlags,
			Runner:          copyGear,
			DatabaseFactory: nilDatabaseFactory,
			Options: map[string]interface{}{
				"allamps":      cpgFlags.Bool("aa", false, "Copy all amps"),
				"allcabs":      cpgFlags.Bool("ac", false, "Copy all cabs"),
				"allfx":        cpgFlags.Bool("ae", false, "Copy all fx"),
				"nocabwithamp": cpgFlags.Bool("c", false, "Don't copy cab with amp"),
				"insertfx":     cpgFlags.Bool("i", false, "Insert fx"),
				"overwritefx":  cpgFlags.Bool("o", false, "Overwrite fx"),
				"recursive":    cpgFlags.Bool("r", false, "Copy to subfolders"),
			},
		},
		"import": {
			Flags:           importFlags,
			Runner:          importPresets,
			DatabaseFactory: secondArgDatabaseFactory,
		},
		"ls": {
			Flags:           lsFlags,
			Runner:          list,
			DatabaseFactory: nilDatabaseFactory,
			Options: map[string]interface{}{
				"fullpath":  lsFlags.Bool("f", false, "Display full path"),
				"recursive": lsFlags.Bool("r", false, "List subfolders"),
			},
		},
		"lsg": {
			Flags:           lsgFlags,
			Runner:          listGear,
			DatabaseFactory: nilDatabaseFactory,
			Options: map[string]interface{}{
				"details": lsgFlags.Bool("d", false, "Show all details"),
				"raw":     lsgFlags.Bool("r", false, "Display raw file"),
			},
		},
		"mkdir": {
			Flags:           mkDirFlags,
			Runner:          makeFolder,
			DatabaseFactory: nilDatabaseFactory,
		},
		"mv": {
			Flags:           mvFlags,
			Runner:          move,
			DatabaseFactory: defaultDatabaseFactory,
		},
		"reindex": {
			Flags:           reindexFlags,
			Runner:          reindex,
			DatabaseFactory: defaultDatabaseFactory,
		},
		"rm": {
			Flags:           rmFlags,
			Runner:          remove,
			DatabaseFactory: defaultDatabaseFactory,
			Options: map[string]interface{}{
				"recursive": rmFlags.Bool("r", false, "Recursive delete"),
			},
		},
		"rmg": {
			Flags:           rmgFlags,
			Runner:          removeGear,
			DatabaseFactory: defaultDatabaseFactory,
			Options: map[string]interface{}{
				"recursive": rmgFlags.Bool("r", false, "Recursive delete gear"),
			},
		},
		"sg": {
			Flags:           sgFlags,
			Runner:          setGear,
			DatabaseFactory: defaultDatabaseFactory,
			Options: map[string]interface{}{
				"recursive": sgFlags.Bool("r", false, "Recursively set gear attribute"),
			},
		},
	}

	command := commands[cmd]

	if command == nil {
		return errors.New("Unknown command " + cmd)
	}

	command.Flags.Parse(args)

	context := ExecutionContext{
		Args:    command.Flags.Args(),
		Options: command.Options,
	}

	database, err := command.DatabaseFactory(context)

	if err != nil {
		return err
	}

	if database == nil {
		return command.Runner(context)
	} else {
		context.Database = database
		return tx(command.Runner, context)
	}

}
