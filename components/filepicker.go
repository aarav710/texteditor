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
	currSearch   string
	displayFiles []string
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
	for _, file := range filePicker.files {
		filePicker.displayFiles = append(filePicker.displayFiles, file.Name())
	}
	if filePicker.currDir != "" {
		filePicker.displayFiles = prepend[string](filePicker.displayFiles, "..")
	}
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
	s.WriteString(fmt.Sprintf("Search: %s\n", m.currSearch))
	for i, file := range m.displayFiles {
		str := fmt.Sprintf("%d. %s\n", i+1, file)
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
		case "ctrl+q":
			return m, tea.Quit
		case "up":
			m.currIndex = max(m.currIndex-1, 0)
		case "down":
			m.currIndex = min(m.currIndex+1, len(m.displayFiles)-1)
		case "enter":
			if m.currIndex == 0 && m.currDir != "" {
				parentDirName := m.parentDir(m.currDir)
				m.switchDir(parentDirName)
			} else {
				fileSelected := m.displayFiles[m.currIndex]
				if m.isDir(fileSelected) {
					if m.currDir == "" {
						m.switchDir(fileSelected)
					} else {
						m.switchDir(m.currDir + "/" + fileSelected)
					}
				} else {
					if getDirectoryName(m.currDir) != "./" {
						fileSelected = getDirectoryName(m.currDir) + "/" + fileSelected
					}
					m.editor.switchFile(fileSelected)
					m.quitting = true
				}
			}
		case "delete", "backspace":
			if m.currSearch != "" {
				m.currSearch = m.currSearch[:len(m.currSearch)-1]
				m.displayFiles = make([]string, 0)
				m.currIndex = 0
				fileNames := make([]string, 0)
				for _, file := range m.files {
					fileNames = append(fileNames, file.Name())
				}
				m.displayFiles = m.FuzzySearch(m.currSearch, fileNames)
				if m.currDir != "" {
					m.displayFiles = prepend[string](m.displayFiles, "..")
				}
			}
		default:
			m.currSearch += msg.String()
			m.currIndex = 0
			m.displayFiles = m.FuzzySearch(m.currSearch, m.displayFiles)
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

func (m *FilePicker) FuzzySearch(search string, fileNames []string) []string {
	results := make([]string, 0)
	for _, fileName := range fileNames {
		i := 0
		j := 0
		for i < len(fileName) && j < len(search) {
			if fileName[i] == search[j] {
				j++
			}
			i++
		}
		if i <= len(fileName) && j == len(search) {
			results = append(results, fileName)
		}
	}
	return results
}

func (m *FilePicker) isDir(name string) bool {
	for _, file := range m.files {
		filename := file.Name()
		if filename == name {
			return file.IsDir()
		}
	}
	return false
}

func (m *FilePicker) switchDir(dirName string) {
	m.currDir = dirName
	files, err := os.ReadDir("./" + m.currDir)
	if err != nil {
		log.Fatal(err)
	}
	m.files = files
	m.displayFiles = make([]string, 0)
	if m.currDir != "" {
		m.displayFiles = append(m.displayFiles, "..")
	}
	m.currIndex = 0
	for _, file := range m.files {
		m.displayFiles = append(m.displayFiles, file.Name())
	}
}

func (m *FilePicker) parentDir(dirName string) string {
	parentDir := ""
	parts := strings.Split(dirName, "/")
	for i := 0; i < len(parts)-1; i++ {
		parentDir += parts[i] + "/"
	}
	if parentDir != "" {
		parentDir = parentDir[:len(parentDir)-1]
	}
	return parentDir
}

func prepend[T any](arr []T, val T) []T {
	result := make([]T, 0)
	result = append(result, val)
	for _, item := range arr {
		result = append(result, item)
	}
	return result
}
