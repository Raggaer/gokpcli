package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sethvargo/go-password/password"
	"github.com/tobischo/gokeepasslib/v2"
	"github.com/tobischo/gokeepasslib/v2/wrappers"
)

const recycleBinGroup = "Recycle Bin"

type newEntryForm struct {
	Title    string
	URL      string
	Password string
	Username string
}

type editEntryForm struct {
	Id       gokeepasslib.UUID
	Title    string
	URL      string
	Password string
	Username string
}

// Command "edit" modifies an entry
func edit(args []string) {
	entry := getEntryByNameOrId(args[0])
	if entry == nil {
		return
	}
	fmt.Print("- Entry username (" + entry.GetContent("UserName") + "): ")
	activeForm = &form{
		Fn: editEntry,
		Data: &editEntryForm{
			Id:       entry.UUID,
			Title:    entry.GetTitle(),
			Username: entry.GetContent("UserName"),
			URL:      entry.GetContent("URL"),
			Password: entry.GetPassword(),
		},
	}
}

func editEntry(f *form, input string) {
	data, ok := f.Data.(*editEntryForm)
	if !ok {
		return
	}
	switch f.Stage {
	case 1:
		if input != "" {
			data.Username = input
		}
		fmt.Print("- Entry title (" + data.Title + "): ")
	case 2:
		if input != "" {
			data.Title = input
		}
		fmt.Print("- Entry password (" + data.Password + "): ")
	case 3:
		if input != "" {
			data.Password = input
		}
		fmt.Print("- Entry URL (" + data.URL + "): ")
	case 4:
		if input != "" {
			data.URL = input
		}
		// Update entry based on UUID
		g := currentGroup()
		for i, e := range g.Entries {
			if e.UUID.Compare(data.Id) {
				fmt.Printf("Entry '%s' modified\r\n", e.GetTitle())
				values := g.Entries[i].Values

				// Update values slice
				for i, v := range values {
					if v.Key == "Title" {
						values[i].Value.Content = data.Title
					}
					if v.Key == "URL" {
						values[i].Value.Content = data.URL
					}
					if v.Key == "UserName" {
						values[i].Value.Content = data.Username
					}
					if v.Key == "Password" {
						values[i].Value.Content = data.Password
					}
				}

				// Update entry modified
				g.Entries[i].Times.LastModificationTime = &wrappers.TimeWrapper{
					Formatted: true,
					Time:      time.Now().In(time.UTC),
				}
				fmt.Print("Database was changed. Save database? (y/N): ")
				activeForm = &form{
					Fn: databaseChangedSaveAlert,
				}
				return
			}
		}
	}
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

func getEntryByPath(path string) *gokeepasslib.Entry {
	// Try to retrieve entry by path
	movedPaths := 0
	paths := strings.Split(path, "/")
pathRange:
	for i, path := range paths {
		// If its the entry from the path
		if i == len(paths)-1 {
			e := getEntryByNameOrId(path)
			// Remove moved paths
			groupHistory = groupHistory[0 : len(groupHistory)-movedPaths]
			return e
		}

		// Move to next group
		gid, err := strconv.Atoi(path)
		if err != nil {
			for x, g := range currentGroup().Groups {
				if strings.ToLower(g.Name) == strings.ToLower(path) {
					groupHistory = append(groupHistory, x)
					continue pathRange
				}
			}
		} else {
			gid--
			for x := range currentGroup().Groups {
				if x == gid {
					groupHistory = append(groupHistory, x)
					movedPaths++
					continue pathRange
				}
			}
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
				entry.Times.LocationChanged = &wrappers.TimeWrapper{
					Formatted: true,
					Time:      time.Now().In(time.UTC),
				}
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
			database.Content.Root.Groups[0].Groups[i].Entries = append(database.Content.Root.Groups[0].Groups[i].Entries, *entry)
			return
		}
	}
	bin := gokeepasslib.NewGroup()
	bin.Name = recycleBinGroup
	bin.Entries = append(bin.Entries, *entry)
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
	entry := getEntryByPath(args[0])
	if entry == nil {
		return
	}
	fmt.Println("Title: " + entry.GetTitle())
	fmt.Println("Uname: " + entry.GetContent("UserName"))
	fmt.Println("Password: " + mask(entry.GetPassword(), "*"))
	fmt.Println("URL: " + entry.GetContent("URL"))
	fmt.Println("Notes:")
	if len(entry.GetContent("Notes")) > 0 {
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
		fmt.Print("- Entry password (use 'gen', 'gen-simple', 'gen-complex' to generate passwords): ")
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
	case "gen-simple":
		return password.Generate(22, 4, 0, false, true)
	case "gen-complex":
		return password.Generate(22, 5, 6, false, false)
	case "gen-large":
		return password.Generate(30, 7, 7, false, false)
	default:
		return input, nil
	}
}

// Command "xp" copies an entry password
func xp(args []string) {
	entry := args[0]
	e := getEntryByPath(entry)
	if e == nil {
		return
	}
	// Update entry access time
	e.Times.LastAccessTime = &wrappers.TimeWrapper{
		Formatted: true,
		Time:      time.Now().In(time.UTC),
	}
	e.Times.UsageCount++
	clipboard.WriteAll(e.GetPassword())
	fmt.Printf("Copied entry '%s' password to clipboard\r\n", e.GetTitle())
	go clipboardClear(clipboardClearDuration, e.GetPassword())
}

// Command "xw" copies an entry URL
func xw(args []string) {
	entry := args[0]
	e := getEntryByPath(entry)
	if e == nil {
		return
	}
	// Update entry access time
	e.Times.LastAccessTime = &wrappers.TimeWrapper{
		Formatted: true,
		Time:      time.Now().In(time.UTC),
	}
	e.Times.UsageCount++
	clipboard.WriteAll(e.GetContent("URL"))
	fmt.Printf("Copied entry '%s' URL to clipboard\r\n", e.GetTitle())
	go clipboardClear(clipboardClearDuration, e.GetContent("URL"))
}

// Command "xu" copies an entry username
func xu(args []string) {
	entry := args[0]
	e := getEntryByPath(entry)
	if e == nil {
		return
	}
	// Update entry access time
	e.Times.LastAccessTime = &wrappers.TimeWrapper{
		Formatted: true,
		Time:      time.Now().In(time.UTC),
	}
	e.Times.UsageCount++
	clipboard.WriteAll(e.GetContent("UserName"))
	fmt.Printf("Copied entry '%s' username to clipboard\r\n", e.GetTitle())
	go clipboardClear(clipboardClearDuration, e.GetContent("UserName"))
}

func mask(data, m string) string {
	return strings.Repeat(m, len(data))
}
