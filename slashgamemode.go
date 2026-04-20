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
