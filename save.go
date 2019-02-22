package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tobischo/gokeepasslib/v2"
)

var confirmDatabaseSave = false

// Shows a database changed message with a save message
func databaseChangedSaveAlert(f *form, answer string) {
	// Remember to close the form
	defer func() {
		activeForm = nil
	}()
	// First we backup the database
	if err := backupDatabase(); err != nil {
		fmt.Println("Unable to backup database:")
		fmt.Println(err.Error())
		return
	}
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

func backupDatabase() error {
	src, err := os.Open(databaseLocation)
	if err != nil {
		return err
	}
	defer src.Close()
	dir, file := filepath.Split(databaseLocation)
	backupName := filepath.Join(dir, time.Now().Format("2006-01-02_15:04:05")+"_"+file)
	dst, err := os.Create(backupName)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}
