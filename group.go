package main

import (
	"fmt"

	"github.com/tobischo/gokeepasslib/v2"
)

type newGroupForm struct {
	Name  string
	Notes string
}

// Command "ng" starts a new group form
func ng() {
	fmt.Print("- Group name: ")
	activeForm = &form{
		Fn:   createNewGroup,
		Data: &newGroupForm{},
	}
}

func createNewGroup(f *form, input string) {
	// Retrieve form data
	data, ok := f.Data.(*newGroupForm)
	if !ok {
		return
	}

	switch f.Stage {
	case 1:
		data.Name = input
		fmt.Print("- Group notes: ")
	case 2:
		data.Notes = input

		// Save new group and close form
		g := gokeepasslib.NewGroup()
		g.Name = data.Name
		g.Notes = data.Notes
		currentGroup().Groups = append(currentGroup().Groups, g)
		fmt.Print("Database was changed. Save database? (y/N): ")
		activeForm = &form{
			Fn: databaseChangedSaveAlert,
		}
	}
}

func currentGroup() *gokeepasslib.Group {
	g := &database.Content.Root.Groups[0]
	for _, h := range groupHistory {
		g = &g.Groups[h]
	}
	return g
}
