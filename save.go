package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tobischo/gokeepasslib"
)

var confirmDatabaseSave = false

// Shows a database changed message with a save message
func databaseChangedSaveAlert(f *form, answer string) {
	// Remember to close the form
	defer func() {
		activeForm = nil
	}()
	if strings.TrimSpace(answer) != "y" {
		return
	}
	if err := saveDatabase(); err != nil {
		fmt.Println("Unable to save database:")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Database saved")
}

func saveDatabase() error {
	database.LockProtectedEntries()
	f, err := os.OpenFile(databaseLocation, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := gokeepasslib.NewEncoder(f)
	if err := encoder.Encode(database); err != nil {
		return err
	}
	return openDatabaseFile(databaseLocation)
}
