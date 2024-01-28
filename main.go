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
		files = append(files, filename)
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("The file with the name %s does not exist\n", filename)
			os.Exit(1)
		}
	}
	// bubbletea tui initialization
	p := tea.NewProgram(components.NewController(files), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
