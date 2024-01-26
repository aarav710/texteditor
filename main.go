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
	filename := os.Args[1]
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("The file with the name %s does not exist\n", filename)
	}

	// bubbletea tui initialization
	p := tea.NewProgram(components.InitialEditorModel(filename), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}