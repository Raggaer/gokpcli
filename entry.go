package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sethvargo/go-password/password"
	"github.com/tobischo/gokeepasslib/v2"
)

const recycleBinGroup = "Recycle Bin"

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
			// Delete entry from the current group
			g.Entries = append(g.Entries[:k], g.Entries[k+1:]...)

			// Add entry to recycle bin
			if g.Name != recycleBinGroup {
				moveEntryToRecycleBin(entry)
			}
			return true
		}
	}
	return false
}

func moveEntryToRecycleBin(entry *gokeepasslib.Entry) {
	// Recycle bin group is at the root, if it does not exist we create it
	for i, g := range database.Content.Root.Groups[0].Groups {
		if g.Name == recycleBinGroup {
			e := gokeepasslib.NewEntry()
			e.Values = append(entry.Values, gokeepasslib.ValueData{Key: "UserName", Value: gokeepasslib.V{Content: entry.GetContent("UserName")}})
			e.Values = append(entry.Values, gokeepasslib.ValueData{Key: "Title", Value: gokeepasslib.V{Content: entry.GetTitle()}})
			e.Values = append(entry.Values, gokeepasslib.ValueData{Key: "URL", Value: gokeepasslib.V{Content: entry.GetContent("URL")}})
			e.Values = append(entry.Values, gokeepasslib.ValueData{Key: "Password", Value: gokeepasslib.V{Content: entry.GetPassword(), Protected: true}})
			database.Content.Root.Groups[0].Groups[i].Entries = append(database.Content.Root.Groups[0].Groups[i].Entries, e)
			return
		}
	}
	bin := gokeepasslib.NewGroup()
	bin.Name = recycleBinGroup
	database.Content.Root.Groups[0].Groups = append(database.Content.Root.Groups[0].Groups, bin)
}

// Command "search" searches for entries of the current group
func search(args []string) {
	// We combine args again to allow spaces
	search := strings.ToLower(strings.Join(args, " "))
	for i, entry := range currentGroup().Entries {
		if !fuzzy.Match(search, strings.ToLower(entry.GetTitle())) {
			continue
		}
		fmt.Printf("%d. %s\r\n", i+1, entry.GetTitle())
	}
}

// Command "show" shows information about an entry
func show(args []string) {
	entry := getEntryByNameOrId(args[0])
	if entry == nil {
		return
	}
	fmt.Println("Title: " + entry.GetTitle())
	fmt.Println("Uname: " + entry.GetContent("UserName"))
	fmt.Println("Password: " + mask(entry.GetPassword(), "*"))
	fmt.Println("URL: " + entry.GetContent("URL"))
	fmt.Println("Notes:")
	if len(entry.GetContent("notes")) > 0 {
		fmt.Println(entry.GetContent("Notes"))
	}
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

// Command "new" starts a new entry form
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

// Command "xw" copies an entry URL
func xp(args []string) {
	entry := args[0]
	e := getEntryByNameOrId(entry)
	if e == nil {
		return
	}
	clipboard.WriteAll(e.GetContent("URL"))
	fmt.Printf("Copied entry '%s' URL to clipboard\r\n", e.GetTitle())
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

func mask(data, m string) string {
	return strings.Repeat(m, len(data))
}
