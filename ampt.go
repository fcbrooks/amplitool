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
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
)

var out io.Writer = os.Stdout

func main() {

	if len(os.Args) < 1 {
		fmt.Fprintln(out, "Subcommand missing")
		os.Exit(1)
	}

	err := ExecuteCommand(os.Args[1], os.Args[2:])

	if err != nil {
		fmt.Fprintln(out, err)
		os.Exit(1)
	}

}
