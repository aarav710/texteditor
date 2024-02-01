package components

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FilePicker struct {
	fileSelected string
	files        []fs.DirEntry
	currDir      string
	quitting     bool
	editor       *EditorModel
	currIndex    int
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func getDirectoryName(name string) string {
	if name == "" || name == " " {
		return "./"
	}
	return name
}

func NewFilePicker(editor *EditorModel) FilePicker {
	filePicker := FilePicker{quitting: true, editor: editor}
	parts := strings.Split(editor.Filename, "/")
	if len(parts) != 1 {
		parts = parts[:len(parts)-1]
		for _, part := range parts {
			filePicker.currDir += part + "/"
		}
		filePicker.currDir = filePicker.currDir[:len(filePicker.currDir)-1]
	}
	files, err := os.ReadDir("./" + filePicker.currDir)
	if err != nil {
		log.Fatal(err)
	}
	filePicker.files = files
	return filePicker
}

func (m *FilePicker) View() string {
	if m.quitting {
		return m.editor.View()
	}
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Current directory: %s\n", getDirectoryName(m.currDir)))
	s.WriteString(fmt.Sprintf("Number of files in this directory: %d\n", len(m.files)))
	s.WriteString("Select a file: \n")
	if getDirectoryName(m.currDir) != "./" {
		s.WriteString(fmt.Sprintf("../\n"))
	}
	for i, file := range m.files {
		str := fmt.Sprintf("%d. %s\n", i+1, file.Name())
		if i == m.currIndex {
			fn := func(s ...string) string {
				return selectedItemStyle.Render("> " + strings.Join(s, " "))
			}
			s.WriteString(fn(str))
		} else {
			fn := itemStyle.Render
			s.WriteString(fn(str))
		}
	}

	return s.String()
}

func (m *FilePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.currIndex = max(m.currIndex-1, 0)
		case "down":
			m.currIndex = min(m.currIndex+1, len(m.files)-1)
		case "enter":
			fileSelected := m.files[m.currIndex].Name()
			if getDirectoryName(m.currDir) != "./" {
				fileSelected = getDirectoryName(m.currDir) + "/" + m.files[m.currIndex].Name()
			}
			m.editor.switchFile(fileSelected)
			m.quitting = true
		}
	}

	return m, nil
}

func (m *FilePicker) getSelectedFile(index int) string {
	return m.files[m.currIndex].Name()
}

func (m *FilePicker) Init() tea.Cmd {
	return nil
}
