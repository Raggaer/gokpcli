package main

import (
	"strings"
)

func handleUserInput(input string) {
	args := strings.Split(input, " ")
	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
	}
	if confirmDatabaseSave {
		databaseChangedSaveAlert(args[0])
		confirmDatabaseSave = false
		return
	}
	switch args[0] {
	case "xp":
		if len(args) > 1 {
			xp(args[1])
		}
	case "xu":
		if len(args) > 1 {
			xu(args[1])
		}
	case "ls":
		ls()
	case "cd":
		if len(args) > 1 {
			cd(args[1])
		}
	case "rm":
		if len(args) > 1 {
			rm(args[1])
		}
	}
}
