package main

import (
	"fmt"
	"strconv"

	"github.com/tobischo/gokeepasslib"
)

var (
	groupHistory []int
)

func currentGroup() *gokeepasslib.Group {
	g := &database.Content.Root.Groups[0]
	for _, h := range groupHistory {
		g = &g.Groups[h]
	}
	return g
}

// Command "ls" shows all groups and entries
// of the current group
func ls() {
	// First show groups
	if len(currentGroup().Groups) > 0 {
		fmt.Println("=== Groups ===")
		for k, g := range currentGroup().Groups {
			fmt.Printf("%d - %s/\r\n", k, g.Name)
		}
	}

	// Show entries
	if len(currentGroup().Entries) > 0 {
		fmt.Println("=== Entries ===")
		for k, e := range currentGroup().Entries {
			fmt.Printf("%d - %s/\r\n", k, e.GetTitle())
		}
	}
}

// Command "cd" changes the current group
func cd(dst string) {
	// Check if going back
	if dst == ".." {
		if len(groupHistory) <= 0 {
			return
		}
		groupHistory = groupHistory[:len(groupHistory)-1]
		if len(groupHistory) > 0 {
			waitCommandMessage = ">> gokpcli/" + currentGroup().Name + " "
		} else {
			waitCommandMessage = ">> gokpcli "
		}
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
				waitCommandMessage = ">> gokpcli/" + currentGroup().Name + " "
				return
			}
		}
		return
	}
	for k := range currentGroup().Groups {
		if k == gid {
			groupHistory = append(groupHistory, k)
			waitCommandMessage = ">> gokpcli/" + currentGroup().Name + " "
			return
		}
	}
}
