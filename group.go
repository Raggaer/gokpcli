package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tobischo/gokeepasslib/v2"
)

type newGroupForm struct {
	Name  string
	Notes string
}

// Command "mkdir" starts a new group form
func mkdir(args []string) {
	fmt.Print("- Group name: ")
	activeForm = &form{
		Fn:   createNewGroup,
		Data: &newGroupForm{},
	}
}

// Command "rmdir" deletes a group from the current working group
func rmdir(args []string) {
	group := getGroupByNameOrId(args[0])
	log.Println(group)
	if group == nil || args[0] == recycleBinGroup {
		return
	}
	if !deleteGroup(group) {
		return
	}
	fmt.Printf("Group '%s' was removed\r\n", group.Name)
	fmt.Print("Database was changed. Save database? (y/N): ")
	activeForm = &form{
		Fn: databaseChangedSaveAlert,
	}
}

func deleteGroup(group *gokeepasslib.Group) bool {
	g := currentGroup()
	for k, e := range g.Groups {
		if strings.ToLower(e.Name) == strings.ToLower(group.Name) {
			// Delete group from the current group
			g.Groups = append(g.Groups[:k], g.Groups[k+1:]...)

			// Move group to recycle bin
			moveGroupToRecycleBin(group)
			return true
		}
	}
	return false
}

func getGroupByNameOrId(group string) *gokeepasslib.Group {
	eid, err := strconv.Atoi(group)
	if err != nil {
		for _, e := range currentGroup().Groups {
			if strings.ToLower(e.Name) == strings.ToLower(group) {
				return &e
			}
		}
		return nil
	}
	eid--
	for k := range currentGroup().Groups {
		if k == eid {
			e := currentGroup().Groups[eid]
			return &e
		}
	}
	return nil
}

func moveGroupToRecycleBin(group *gokeepasslib.Group) {
	// Recycle bin group is at the root, if it does not exist we create it
	for i, g := range database.Content.Root.Groups[0].Groups {
		if g.Name == recycleBinGroup {
			n := gokeepasslib.NewGroup()
			n.Name = group.Name
			n.Entries = group.Entries
			database.Content.Root.Groups[0].Groups[i].Groups = append(database.Content.Root.Groups[0].Groups[i].Groups, n)
			return
		}
	}
	bin := gokeepasslib.NewGroup()
	bin.Name = recycleBinGroup
	database.Content.Root.Groups[0].Groups = append(database.Content.Root.Groups[0].Groups, bin)

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

func getGroupByPath(path string) *gokeepasslib.Group {
	// Try to retrieve entry by path
	movedPaths := 0
	paths := strings.Split(path, "/")
pathRange:
	for i, path := range paths {
		// If its the entry from the path
		if i == len(paths)-1 {
			e := currentGroup()
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

func currentGroup() *gokeepasslib.Group {
	g := &database.Content.Root.Groups[0]
	for _, h := range groupHistory {
		g = &g.Groups[h]
	}
	return g
}
