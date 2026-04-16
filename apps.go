package main

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type App struct {
	Name string
	Exec string
}

func getApps() []App {
	var apps []App
	dirs := []string{"/usr/share/applications", os.Getenv("HOME") + "/.local/share/applications"}
	for _, dir := range dirs {
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if strings.HasSuffix(path, ".desktop") {
				if name, exec := parseDesktop(path); name != "" {
					apps = append(apps, App{Name: name, Exec: exec})
				}
			}
			return nil
		})
	}
	return apps
}

func parseDesktop(path string) (name, exec string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	isDesktopEntry := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "[Desktop Entry]" {
			isDesktopEntry = true
		} else if isDesktopEntry && strings.HasPrefix(line, "Name=") {
			name = strings.TrimPrefix(line, "Name=")
		} else if isDesktopEntry && strings.HasPrefix(line, "Exec=") {
			exec = strings.TrimPrefix(line, "Exec=")
			// Strip field codes like %u, %f, etc.
			exec = strings.Split(exec, " %")[0]
		}
		if line != "" && strings.HasPrefix(line, "[") && line != "[Desktop Entry]" {
			isDesktopEntry = false
		}
	}
	return
}