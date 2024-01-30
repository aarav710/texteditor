package components

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"texteditor/textctrl"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type EditorModel struct {
	CursorPositionX int
	CursorPositionY int
	LinesCount      int
	Content         []textinput.Model
	Filename        string
	IsInsertMode    bool
	cursor          cursor.Mode
	textctrl        *textctrl.Handler
	linesDisplayed  int
	fileSelector    *FilePicker
}

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Background(lipgloss.Color("#3C3C3C"))
	footerStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Background(lipgloss.Color("#3C3C3C"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

// try finding a way to not assign default height but actually calculate the size
const defaultLinesDisplayed = 31

func NewTextInput(content string, lineNo int) textinput.Model {
	textinput := textinput.New()
	textinput.SetValue(content)
	textinput.Cursor.Style = cursorStyle
	w, h, _ := term.GetSize(0)
	textinput.Width = w
	textinput.TextStyle.Height(h / defaultLinesDisplayed)
	if lineNo == 0 {
		textinput.SetCursor(0)
		textinput.TextStyle = focusedStyle
		textinput.Focus()
	} else {
		textinput.TextStyle = blurredStyle
		textinput.Blur()
	}
	return textinput
}

func readContent(filename string) ([]textinput.Model, int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("The file with the name %s could not be opened\n", filename)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	content := make([]textinput.Model, 0)
	linesCount := 0
	for scanner.Scan() {
		content = append(content, NewTextInput(scanner.Text(), linesCount))
		linesCount++
	}
	return content, linesCount
}

func (m *EditorModel) switchFile(filename string) {
	m.Filename = filename
	m.Content, m.LinesCount = readContent(m.Filename)
}

func InitialEditorModel(filename string) *EditorModel {
	m := EditorModel{}
	m.Filename = filename
	m.Content, m.LinesCount = readContent(m.Filename)
	m.textctrl = textctrl.NewHandler()
	m.linesDisplayed = defaultLinesDisplayed
	fileSelecter := NewFilePicker(&m)
	m.fileSelector = &fileSelecter
	return &m
}

func LineView(row int, m *EditorModel, linesCountLength int) string {
	var b strings.Builder
	lineNo := m.CursorPositionY - row
	if row == m.CursorPositionY {
		lineNo = m.CursorPositionY + 1
	} else if lineNo < 0 {
		lineNo *= -1
	}
	lineSpacing := fmt.Sprintf("%d", lineNo)
	digitLen := digitLength(lineNo)
	for j := 0; j < linesCountLength-digitLen+2; j++ {
		lineSpacing += " "
	}
	b.WriteString(lineSpacing)
	b.WriteString(m.Content[row].View())
	if row < m.LinesCount-1 {
		b.WriteRune('\n')
	}
	return b.String()
}

func (m *EditorModel) View() string {
	if !m.fileSelector.quitting {
		return m.fileSelector.View()
	}
	var upperStr string
	var lowerStr string
	LinesCountLength := digitLength(m.LinesCount)

	left, right := m.CursorPositionY, m.CursorPositionY+1
	total := 0
	for (left >= 0 || right < m.LinesCount) && total < m.linesDisplayed {
		if left >= 0 {
			lowerStr = LineView(left, m, LinesCountLength) + lowerStr
			total++
			left--
		}
		if total >= m.linesDisplayed {
			break
		}
		if right < m.LinesCount {
			upperStr += LineView(right, m, LinesCountLength)
			total++
			right++
		}
	}
	footer := textinput.New()
	footerValue := m.Filename
	w, _, _ := term.GetSize(0)
	for i := 40; i < w; i++ {
		footerValue += " "
	}
	percentDone := int((float32(m.CursorPositionY) / float32(m.LinesCount)) * 100)
	xCursorValue := m.CursorPositionX
	if xCursorValue == math.MaxInt {
		xCursorValue = len(m.Content[m.CursorPositionY].Value())
	}
	footerValue += fmt.Sprintf("%d,%d,%d%%", m.CursorPositionY+1, xCursorValue, percentDone)
	footer.SetValue(footerValue)
	footer.Blur()
	footer.Width = w - 40
	footer.TextStyle = footerStyle
	return lowerStr + upperStr + fmt.Sprintf("\n") + footer.View()
}

func (m *EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			if m.CursorPositionX == math.MaxInt {
				m.CursorPositionX = len(m.Content[m.CursorPositionY].Value()) - 1
			} else if m.CursorPositionX > 0 {
				m.CursorPositionX--
			}
			m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
		case "j":
			if m.CursorPositionY < m.LinesCount-1 {
				switchBlurToFocus(&m.Content[m.CursorPositionY], &m.Content[m.CursorPositionY+1])
				m.CursorPositionY++
				m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
			}
		case "k":
			if m.CursorPositionY > 0 {
				switchBlurToFocus(&m.Content[m.CursorPositionY], &m.Content[m.CursorPositionY-1])
				m.CursorPositionY--
				m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
			}
		case "l":
			if m.CursorPositionX < len(m.Content[m.CursorPositionY].Value()) {
				m.CursorPositionX++
				m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
			}
		case "i", "a":
			m.IsInsertMode = true
			if msg.String() == "a" {
				m.CursorPositionX++
			}
		case "o":
			m.IsInsertMode = true
			m.LinesCount++
			m.CursorPositionY++
		case "G":
			newYPosition := m.LinesCount - 1
			switchBlurToFocus(&m.Content[m.CursorPositionY], &m.Content[newYPosition])
			m.CursorPositionY = newYPosition
			m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
		case "esc", "ctrl+c":
			if !m.fileSelector.quitting {
				m.fileSelector.quitting = true
			} else {
				m.IsInsertMode = false
			}
		case "ctrl+d":
			newYPosition := min(m.LinesCount-1, m.CursorPositionY+15)
			switchBlurToFocus(&m.Content[m.CursorPositionY], &m.Content[newYPosition])
			m.CursorPositionY = newYPosition
			m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
		case "ctrl+u":
			newYPosition := max(0, m.CursorPositionY-15)
			switchBlurToFocus(&m.Content[m.CursorPositionY], &m.Content[newYPosition])
			m.CursorPositionY = newYPosition
			m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
		case "0":
			m.CursorPositionX = 0
			m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
		case "$":
			m.CursorPositionX = math.MaxInt
			m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
		case "^":
			// todo
			str := m.Content[m.CursorPositionY].Value()
			width := len(str)
			newCursorPositionX := 0
			for string(str[newCursorPositionX]) == " " && newCursorPositionX < width {
				newCursorPositionX++
			}
			m.CursorPositionX = newCursorPositionX
			m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
		case "ctrl+f":
			m.fileSelector.quitting = false
		case "down", "up", "enter":
			if !m.fileSelector.quitting {
				m.fileSelector.Update(msg)
			}
		default:
			if m.textctrl.IsValidMotion() {
				m.textctrl.ExecuteMotion()
			} else {
				m.textctrl.AddToCurrMotion(msg.String())
			}
		}
	}
	return m, nil
}

func (m *EditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func digitLength(digit int) int {
	length := 1
	for digit > 0 {
		length += 1
		digit /= 10
	}
	return length
}

func switchBlurToFocus(old *textinput.Model, newFocus *textinput.Model) {
	old.Blur()
	old.TextStyle = blurredStyle
	old.Cursor.TextStyle = blurredStyle
	newFocus.Focus()
	newFocus.TextStyle = focusedStyle
	newFocus.Cursor.TextStyle = focusedStyle
}
