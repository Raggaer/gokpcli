package main

import (
	"fmt"
	"log"
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
		Key:       "exit",
		Fn:        exit,
		Help:      "Exits the application",
		HelpSmall: "Exits the application",
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
		Key:       "ng",
		Fn:        ng,
		Help:      "Creates a new group inside the current group",
		HelpSmall: "Creates a new group inside the current group",
	},
	{
		Key:       "ne",
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
			log.Println(args)
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
	for _, command := range commands {
		if command.Help == "" && command.HelpSmall == "" {
			continue
		}
		fmt.Println(command.Key + " -- " + command.HelpSmall)
	}
}
