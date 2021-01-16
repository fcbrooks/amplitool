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
	"strings"
)

func openDatabase(dbFile string) (*sql.DB, error) {
	return sql.Open("sqlite3", strings.ReplaceAll(dbFile, "\\", "/"))
}

func tx(runner Runner, context ExecutionContext) error {

	database := context.Database

	tx, err := database.Begin()

	if err != nil {
		_ = database.Close()
		return err
	}

	rtrn := runner(context)

	if rtrn != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	_ = database.Close()

	return rtrn
}
