package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"texteditor/components"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	filesCount := len(os.Args) - 1
	if filesCount == 0 {
		fmt.Printf("Please pass in files to open in the editor\n")
		os.Exit(1)
	}
	if filesCount > 5 {
		fmt.Printf("Please pass in a maximum of 5 files to be opened\n")
		os.Exit(1)
	}

	files := make([]string, 0)
	for i := 0; i < filesCount; i++ {
		filename := os.Args[i+1]
		if len(filename) >= 3 && string(filename[0]) == "." && string(filename[1]) == "/" {
			filename = filename[2:]
		}
		files = append(files, filename)
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("The file with the name %s does not exist\n", filename)
			os.Exit(1)
		}
	}
	// bubbletea tui initialization
	controller := components.NewController(files)
	p := tea.NewProgram(&controller, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
