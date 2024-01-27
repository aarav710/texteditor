package components

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"texteditor/textctrl"
	"time"

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

const debounceDur = time.Second * 2
const useHighPerfRender = true

// try finding a way to not assign default height but actually calculate the size
const defaultLinesDisplayed = 31

type footer struct {
	filename string
}

func NewTextInput(content string, editor *EditorModel, lineNo int) textinput.Model {
	textinput := textinput.New()
	textinput.SetValue(content)
	textinput.Cursor.Style = cursorStyle
	if editor.CursorPositionY == lineNo {
		textinput.SetCursor(editor.CursorPositionX)
		textinput.TextStyle = focusedStyle
		textinput.Focus()
	} else {
		textinput.TextStyle = blurredStyle
		textinput.Blur()
	}
	return textinput
}

func InitialEditorModel(filename string) EditorModel {
	m := EditorModel{}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("The file with the name %s could not be opened\n", filename)
	}
	defer file.Close()
	m.Filename = filename
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m.Content = append(m.Content, NewTextInput(scanner.Text(), &m, m.LinesCount))
		m.LinesCount++
	}
	m.textctrl = textctrl.NewHandler()
	m.linesDisplayed = defaultLinesDisplayed
	return m
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

func (m EditorModel) View() string {
	var upperStr string
	var lowerStr string
	LinesCountLength := digitLength(m.LinesCount)

	left, right := m.CursorPositionY, m.CursorPositionY+1
	total := 0
	for (left >= 0 || right < m.LinesCount) && total < m.linesDisplayed {
		if left >= 0 {
			lowerStr = LineView(left, &m, LinesCountLength) + lowerStr
			total++
			left--
		}
		if total >= m.linesDisplayed {
			break
		}
		if right < m.LinesCount {
			upperStr += LineView(right, &m, LinesCountLength)
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
	footer.TextStyle = footerStyle
	return lowerStr + upperStr + fmt.Sprintf("\n") + footer.View()
}

func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
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
			m.Content = append(m.Content, NewTextInput("", &m, m.LinesCount))
		case "G":
			newYPosition := m.LinesCount - 1
			switchBlurToFocus(&m.Content[m.CursorPositionY], &m.Content[newYPosition])
			m.CursorPositionY = newYPosition
			m.Content[m.CursorPositionY].SetCursor(m.CursorPositionX)
		case "esc", "ctrl+c":
			m.IsInsertMode = false
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

func (m EditorModel) Init() tea.Cmd {
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
