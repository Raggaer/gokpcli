package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"
)

type shortcut struct {
	Alias   string
	Command string
	Input   []string
}

func loadShortcuts(path string) []*shortcut {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(content))
	shortcuts := []*shortcut{}
	for scanner.Scan() {
		tt := strings.TrimSpace(scanner.Text())
		// Skip comments
		if strings.HasPrefix(tt, "#") {
			continue
		}
		parts := strings.Split(tt, " ")
		if len(parts) <= 1 {
			continue
		}
		alias, parts := parts[0], parts[1:]
		command, parts := parts[0], parts[1:]
		s := &shortcut{
			Command: command,
			Alias:   alias,
			Input:   parts,
		}
		shortcuts = append(shortcuts, s)
	}
	return shortcuts
}
