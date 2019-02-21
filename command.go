package main

import (
	"fmt"
	"strings"
)

type command struct {
	Key       string
	Fn        func(args []string)
	Help      string
	HelpSmall string
}

var commands = []command{
	{
		Key: "save",
		Fn: func(args []string) {
			if err := saveDatabase(); err != nil {
				fmt.Println("Unable to save database")
				fmt.Println(err.Error())
			}
		},
		Help:      "Saves the database to disk",
		HelpSmall: "Saves the database to disk",
	},
	{
		Key:       "exit",
		Fn:        exit,
		Help:      "Exits the application",
		HelpSmall: "Exits the application",
	},
	{
		Key: "show",
		Fn: func(args []string) {
			if len(args) >= 1 {
				show(args)
			}
		},
		Help:      "Shows an entry (show <entry name|number>)",
		HelpSmall: "Shows an entry (show <entry name|number>)",
	},
	{
		Key: "xp",
		Fn: func(args []string) {
			if len(args) >= 1 {
				xp(args)
			}
		},
		Help:      "Copies an entry password (xp <entry name|number>)",
		HelpSmall: "Copies an entry password (xp <entry name|number>)",
	},
	{
		Key: "xw",
		Fn: func(args []string) {
			if len(args) >= 1 {
				xw(args)
			}
		},
		Help:      "Copies an entry URL (xw <entry name|number>)",
		HelpSmall: "Copies an entry URL (xw <entry name|number>)",
	},
	{
		Key: "xu",
		Fn: func(args []string) {
			if len(args) >= 1 {
				xu(args)
			}
		},
		Help:      "Copies an entry username (xu <entry name|number>)",
		HelpSmall: "Copies an entry username (xu <entry name|number>)",
	},
	{
		Key:       "mkdir",
		Fn:        mkdir,
		Help:      "Creates a new group inside the current group",
		HelpSmall: "Creates a new group inside the current group",
	},
	{
		Key:       "new",
		Fn:        ne,
		Help:      "Creates a new entry inside the current group",
		HelpSmall: "Creates a new entry inside the current group",
	},
	{
		Key:       "ls",
		Fn:        ls,
		Help:      "Lists entries and groups of the current group",
		HelpSmall: "Lists entries and groups of the current group",
	},
	{
		Key: "cd",
		Fn: func(args []string) {
			if len(args) >= 1 {
				cd(args)
			}
		},
		Help:      "Change directory (path to a group)",
		HelpSmall: "Change directory (path to a group)",
	},
	{
		Key: "rm",
		Fn: func(args []string) {
			if len(args) >= 1 {
				rm(args)
			}
		},
		Help:      "Removes an entry of the current group (rm <entry path|number>)",
		HelpSmall: "Removes an entry of the current group (rm <entry path|number>)",
	},
	{
		Key: "search",
		Fn: func(args []string) {
			if len(args) >= 1 {
				search(args)
			}
		},
		Help:      "Performs a fuzzy search on all the current group entries, by title (search <query>)",
		HelpSmall: "Performs a fuzzy search on all the current group entries, by title (search <query>)",
	},
	{
		Key: "rmdir",
		Fn: func(args []string) {
			if len(args) >= 1 {
				rmdir(args)
			}
		},
		Help:      "Deletes a group (rmdir <group_name|number>)",
		HelpSmall: "Deletes a group (rmdir <group_name|number>)",
	},
	{
		Key: "edit",
		Fn: func(args []string) {
			if len(args) >= 1 {
				edit(args)
			}
		},
		Help:      "Modifies an entry (edit <path|number>)",
		HelpSmall: "Modifies an entry (edit <path|number>)",
	},
}

func handleUserInput(input string) {
	args := strings.Split(input, " ")
	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
	}

	// Check if any form is active
	if activeForm != nil {
		activeForm.Stage++
		activeForm.Fn(activeForm, args[0])
		return
	}

	// Execute command
	if args[0] == "help" {
		help(args)
		return
	}
	for _, command := range commands {
		if command.Key == args[0] {
			command.Fn(args[1:])
			break
		}
	}
}

// Command "exit" closes the application
func exit(args []string) {
	close(quit)
}

// Command "help" shows the list of commands
func help(args []string) {
	// Check if we want help about a command
	if len(args) > 1 {
		for _, command := range commands {
			if command.Key == args[1] {
				fmt.Println(command.Key + " -- " + command.Help)
			}
		}
		return
	}
	for _, command := range commands {
		if command.Help == "" && command.HelpSmall == "" {
			continue
		}
		fmt.Println(command.Key + " -- " + command.HelpSmall)
	}
	fmt.Println("\r\nType \"help <command>\" for a more detailed help on a command")
}
