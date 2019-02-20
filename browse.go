package main

import (
	"fmt"
	"strconv"
)

var (
	groupHistory []int
)

// Command "ls" shows all groups and entries of the current group
func ls(args []string) {
	// First show groups
	if len(currentGroup().Groups) > 0 {
		fmt.Println("=== Groups ===")
		for k, g := range currentGroup().Groups {
			fmt.Printf("%d - %s/\r\n", k+1, g.Name)
		}
	}

	// Show entries
	if len(currentGroup().Entries) > 0 {
		fmt.Println("=== Entries ===")
		for k, e := range currentGroup().Entries {
			fmt.Printf("%d - %s/\r\n", k+1, e.GetTitle())
		}
	}
}

// Command "cd" changes the current group
func cd(args []string) {
	dst := args[0]
	// Check if going back
	if dst == ".." {
		if len(groupHistory) <= 0 {
			return
		}
		groupHistory = groupHistory[:len(groupHistory)-1]
		waitCommandMessage = buildApplicationWaitMessage()
		return
	}

	if len(currentGroup().Groups) <= 0 {
		return
	}

	gid, err := strconv.Atoi(dst)
	if err != nil {
		for k, g := range currentGroup().Groups {
			if g.Name == dst {
				groupHistory = append(groupHistory, k)
				waitCommandMessage = buildApplicationWaitMessage()
				return
			}
		}
		return
	}
	for k := range currentGroup().Groups {
		if k == gid-1 {
			groupHistory = append(groupHistory, k)
			waitCommandMessage = buildApplicationWaitMessage()
			return
		}
	}
}
