package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/tobischo/gokeepasslib"
)

func getEntryByNameOrId(entry string) *gokeepasslib.Entry {
	eid, err := strconv.Atoi(entry)
	if err != nil {
		for _, e := range currentGroup().Entries {
			if strings.ToLower(e.GetTitle()) == strings.ToLower(entry) {
				return &e
			}
		}
		return nil
	}
	eid--
	for k := range currentGroup().Entries {
		if k == eid {
			e := currentGroup().Entries[eid]
			return &e
		}
	}
	return nil
}

func deleteEntry(entry *gokeepasslib.Entry) bool {
	g := currentGroup()
	for k, e := range g.Entries {
		if e.UUID.Compare(entry.UUID) {
			g.Entries = append(g.Entries[:k], g.Entries[k+1:]...)
			return true
		}
	}
	return false
}

// Command "rm" removes an entry
func rm(entry string) {
	e := getEntryByNameOrId(entry)
	if e == nil {
		return
	}
	if !deleteEntry(e) {
		return
	}
	fmt.Printf("Entry '%s' removed\r\n", e.GetTitle())
	fmt.Print("Database was changed. Save database? (y/N): ")
	confirmDatabaseSave = true
}

// Command "xp" copies an entry password
func xp(entry string) {
	e := getEntryByNameOrId(entry)
	if e == nil {
		return
	}
	clipboard.WriteAll(e.GetPassword())
	fmt.Printf("Copied entry '%s' password to clipboard\r\n", e.GetTitle())
}

// Command "xu" copies an entry username
func xu(entry string) {
	e := getEntryByNameOrId(entry)
	if e == nil {
		return
	}
	clipboard.WriteAll(e.GetContent("UserName"))
	fmt.Printf("Copied entry '%s' username to clipboard\r\n", e.GetTitle())
}
