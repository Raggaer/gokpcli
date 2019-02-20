package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/sethvargo/go-password/password"
	"github.com/tobischo/gokeepasslib/v2"
)

type newEntryForm struct {
	Title    string
	URL      string
	Password string
	Username string
}

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
func rm(args []string) {
	entry := args[0]
	e := getEntryByNameOrId(entry)
	if e == nil {
		return
	}
	if !deleteEntry(e) {
		return
	}
	fmt.Printf("Entry '%s' removed\r\n", e.GetTitle())
	fmt.Print("Database was changed. Save database? (y/N): ")
	activeForm = &form{
		Fn: databaseChangedSaveAlert,
	}
}

// Command "ne" starts a new entry form
func ne(args []string) {
	fmt.Print("- Entry username: ")
	activeForm = &form{
		Fn:   createNewEntry,
		Data: &newEntryForm{},
	}
}

func createNewEntry(f *form, input string) {
	// Retrieve form data
	data, ok := f.Data.(*newEntryForm)
	if !ok {
		return
	}

	switch f.Stage {
	case 1:
		data.Username = input
		fmt.Print("- Entry title: ")
	case 2:
		data.Title = input
		fmt.Print("- Entry password: ")
	case 3:
		// Here we can generate password
		pw, err := generateEntryPassword(input)
		if err != nil {
			fmt.Println("Unable to generate entry password")
			fmt.Println(err.Error())
			activeForm = nil
			return
		}
		data.Password = pw
		fmt.Print("- Entry url: ")
	case 4:
		data.URL = input

		// Create new entry
		entry := gokeepasslib.NewEntry()
		entry.Values = append(entry.Values, gokeepasslib.ValueData{Key: "UserName", Value: gokeepasslib.V{Content: data.Username}})
		entry.Values = append(entry.Values, gokeepasslib.ValueData{Key: "Title", Value: gokeepasslib.V{Content: data.Title}})
		entry.Values = append(entry.Values, gokeepasslib.ValueData{Key: "URL", Value: gokeepasslib.V{Content: data.URL}})
		entry.Values = append(entry.Values, gokeepasslib.ValueData{Key: "Password", Value: gokeepasslib.V{Content: data.Password, Protected: true}})
		currentGroup().Entries = append(currentGroup().Entries, entry)
		fmt.Print("Database was changed. Save database? (y/N): ")
		activeForm = &form{
			Fn: databaseChangedSaveAlert,
		}
	}
}

func generateEntryPassword(input string) (string, error) {
	switch input {
	// Simple password generation
	case "gen":
		return password.Generate(22, 4, 4, false, true)
	default:
		return input, nil
	}
}

// Command "xp" copies an entry password
func xp(args []string) {
	entry := args[0]
	e := getEntryByNameOrId(entry)
	if e == nil {
		return
	}
	clipboard.WriteAll(e.GetPassword())
	fmt.Printf("Copied entry '%s' password to clipboard\r\n", e.GetTitle())
}

// Command "xu" copies an entry username
func xu(args []string) {
	entry := args[0]
	e := getEntryByNameOrId(entry)
	if e == nil {
		return
	}
	clipboard.WriteAll(e.GetContent("UserName"))
	fmt.Printf("Copied entry '%s' username to clipboard\r\n", e.GetTitle())
}
