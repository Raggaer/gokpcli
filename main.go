package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tobischo/gokeepasslib"
)

var (
	allowMainInput       = true
	quit                 = make(chan struct{})
	msg                  = make(chan string, 1)
	waitCommandMessage   = ">> gokpcli "
	databaseLocation     string
	passwordFileLocation string
	database             *gokeepasslib.Database
)

func main() {
	fmt.Print(waitCommandMessage)

	// Parse application flags
	flag.StringVar(&databaseLocation, "db", "", "KeePass2 database file location")
	flag.StringVar(&passwordFileLocation, "pwfile", "", "File that stores your KeePass2 database")
	flag.Parse()

	// Open database file
	if err := openDatabaseFile(databaseLocation); err != nil {
		fmt.Println("Unable to open KeePass2 database file")
		return
	}

	// Register channel to get signal notifications
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go readUserInput()

	// Read user input until we get a signal
loop:
	for {
		select {
		case <-sigs:
			break loop
		case <-quit:
			break loop
		case s := <-msg:
			handleUserInput(s)
			if !confirmDatabaseSave {
				fmt.Print(waitCommandMessage)
			}
		}
	}
}

func openDatabaseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	database = gokeepasslib.NewDatabase()

	// Read password from file
	pw, err := readPasswordFile(passwordFileLocation)
	if err != nil {
		return err
	}

	database.Credentials = gokeepasslib.NewPasswordCredentials(pw)
	if err := gokeepasslib.NewDecoder(file).Decode(database); err != nil {
		return err
	}
	database.UnlockProtectedEntries()

	return nil
}

func readPasswordFile(path string) (string, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(f)), nil
}

func readUserInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		msg <- text
	}
}
