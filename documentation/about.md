## Slogan

Your right hand

## Description

A fast and efficient application launcher built with Go and Fyne, designed to help you quickly find and launch installed applications on your Linux system.

## Tech Stack

- **Language**: Go (Golang)
- **GUI Framework**: Fyne v2.7.3
- **Platform**: Linux

## Packages Used

### Core Go Packages

- `bufio` - Buffered I/O for reading .desktop files
- `io/fs` - Filesystem interfaces for directory traversal
- `os` - Operating system interfaces (environment variables, file operations)
- `os/exec` - Executing external commands to launch applications
- `path/filepath` - Path manipulation for cross-platform compatibility
- `strings` - String manipulation for parsing and searching

### Fyne GUI Packages

- `fyne.io/fyne/v2` - Core Fyne framework
- `fyne.io/fyne/v2/app` - Application management
- `fyne.io/fyne/v2/container` - Layout containers (e.g., Border layout)
- `fyne.io/fyne/v2/widget` - UI widgets (Entry, List, etc.)

## Features

- [ ] Search the web
  - [ ] Site search
- [ ] Search local files
- [ ] Search local folders
- [ ] Offline translation (ar and en)
- [ ] Online translation (ar and en)
- [x] App launcher
- [ ] Script runner
- [ ] Offline ticktick
- [ ] Ticktick integration
