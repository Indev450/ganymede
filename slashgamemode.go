/*
	/gamemode command implementation
	Copyright (C) 2026 Indev

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
	"os"
	"fmt"
	"bufio"
	"strings"
)

func doSlashGamemode(path string) string {
	file, err := os.Open(path)

	if err != nil {
		fmt.Printf("Failed to open %s: %s\n", path, err)
		return "Unable to fetch gamemodes right now."
	}

	defer file.Close()

	var list []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		list = append(list, line)
	}

	if len(list) == 0 {
		list = append(list, "regular")
	}

	if len(list) == 1 {
		return fmt.Sprintf("Current gamemode is: %s.", list[0])
	} else {
		var builder strings.Builder

		fmt.Fprintf(&builder, "Current gamemodes are: ")

		for i, name := range list {
			fmt.Fprintf(&builder, "%s", name)

			if i < len(list)-1 {
				fmt.Fprintf(&builder, ", ")
			}
		}

		return builder.String()
	}
}
