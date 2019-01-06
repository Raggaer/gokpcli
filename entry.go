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
	for k := range currentGroup().Entries {
		if k == eid {
			e := currentGroup().Entries[eid]
			return &e
		}
	}
	return nil
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
