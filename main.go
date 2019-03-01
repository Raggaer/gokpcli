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
	"time"

	"github.com/atotto/clipboard"
	"github.com/tobischo/gokeepasslib/v2"
)

var (
	clipboardClearDuration = time.Second * 10
	quit                   = make(chan struct{})
	msg                    = make(chan string, 1)
	waitCommandMessage     = ">> gokpcli "
	databaseLocation       string
	passwordFileLocation   string
	database               *gokeepasslib.Database
	doNotBackups           = false
)

func main() {
	// Parse application flags
	flag.StringVar(&databaseLocation, "db", "", "KeePass2 database file location")
	flag.StringVar(&passwordFileLocation, "pwfile", "", "File that stores your KeePass2 database password")
	flag.BoolVar(&doNotBackups, "nbackup", false, "Do not use the builtin backup method")
	flag.Parse()

	// Open database file
	if err := openDatabaseFile(databaseLocation); err != nil {
		fmt.Println("Unable to open KeePass2 database file")
		fmt.Println(err.Error())
		return
	}

	fmt.Println("gokpcli is ready for operation")
	fmt.Println("Type 'help' for a description of available commands")
	fmt.Println("Type 'help <command>' for details on individual commands")
	fmt.Print(waitCommandMessage)

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
			if activeForm == nil {
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

func buildApplicationWaitMessage() string {
	g := database.Content.Root.Groups[0]
	str := ">> gokpcli"
	for i, h := range groupHistory {
		if i == 0 {
			str += "/"
		}
		g = g.Groups[h]
		str += g.Name + "/"
	}
	str += " "
	return str

}

// Clears the clipboard only if the content is equal to clipboard
func clipboardClear(at time.Duration, content string) {
	time.Sleep(at)
	clip, err := clipboard.ReadAll()
	if err != nil {
		return
	}
	if clip == content {
		clipboard.WriteAll("")
	}
}
