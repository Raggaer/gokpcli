package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

type command struct {
	Key       string
	Fn        func(args []string)
	Help      string
	HelpSmall string
}

var (
	shortcuts = []*shortcut{}
	commands  = []command{
		{
			Key: "backup",
			Fn: func(args []string) {
				if err := backupDatabase(); err != nil {
					fmt.Println("Unable to backup database:")
					fmt.Println(err.Error())
				} else {
					fmt.Println("Database backup created")
				}
			},
			Help:      "Backups the database file. The new backup is saved as yyyy-mm-dd_hh-ii-ss_name.kdbx",
			HelpSmall: "Backups the database file. The new backup is saved as yyyy-mm-dd_hh-ii-ss_name.kdbx",
		},
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
			Help:      "Shows an entry (show <entry name|number|entry path>)",
			HelpSmall: "Shows an entry (show <entry name|number|entry path>)",
		},
		{
			Key: "xp",
			Fn: func(args []string) {
				if len(args) >= 1 {
					xp(args)
				}
			},
			Help:      "Copies an entry password (xp <entry name|number|entry path>)",
			HelpSmall: "Copies an entry password (xp <entry name|number|entry path>)",
		},
		{
			Key: "xw",
			Fn: func(args []string) {
				if len(args) >= 1 {
					xw(args)
				}
			},
			Help:      "Copies an entry URL (xw <entry name|number|entry path>)",
			HelpSmall: "Copies an entry URL (xw <entry name|number|entry path>)",
		},
		{
			Key: "xu",
			Fn: func(args []string) {
				if len(args) >= 1 {
					xu(args)
				}
			},
			Help:      "Copies an entry username (xu <entry name|number|entry path>)",
			HelpSmall: "Copies an entry username (xu <entry name|number|entry path>)",
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
			Help:      "Removes an entry of the current group (rm <entry name|number>)",
			HelpSmall: "Removes an entry of the current group (rm <entry name|number>)",
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
			Help:      "Deletes a group (rmdir <group name|number>)",
			HelpSmall: "Deletes a group (rmdir <group name|number>)",
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
		{
			Key:       "xx",
			Fn:        xx,
			Help:      "Clears the clipboard",
			HelpSmall: "Clears the clipboard",
		},
	}
)

func handleUserInput(input string) {
	args := strings.Split(input, " ")
	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
	}

	// Check if any form is active
	if activeForm != nil {
		activeForm.Stage++
		activeForm.Fn(activeForm, strings.TrimSpace(input))
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
			return
		}
	}

	// If no command match then try to match against shortcuts
	for _, shortcut := range shortcuts {
		if shortcut.Alias == args[0] {
			for _, command := range commands {
				if command.Key == shortcut.Command {
					command.Fn(shortcut.Input)
					return
				}
			}
		}
	}
}

// Command "exit" closes the application
func exit(args []string) {
	close(quit)
}

// Command "xx" clears the clipboard
func xx(args []string) {
	clipboard.WriteAll("")
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
	fmt.Println("\r\nType \"help <command>\" for a more detailed help on a command\r\n")
	for _, shortcut := range shortcuts {
		fmt.Println("Shortcut '" + shortcut.Alias + "' -- " + shortcut.Command + " " + strings.Join(shortcut.Input, " "))
	}
}
