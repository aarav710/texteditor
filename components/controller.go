package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type Controller struct {
	activeTab int
	tabs      []*EditorModel
	files     []string
}

const maxTabs = 5

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "86", Dark: "86"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func NewController(files []string) Controller {
	controller := Controller{files: files, activeTab: 0}
	controller.tabs = make([]*EditorModel, len(files))
	for i := 0; i < len(files); i++ {
		controller.tabs[i] = InitialEditorModel(files[i])
	}
	return controller
}

func (m *Controller) Init() tea.Cmd {
	return nil
}

func (m *Controller) View() string {
	doc := strings.Builder{}
	var renderedTabs []string

	for i := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		displayName := m.tabs[i].Filename
		if len(displayName) >= 2 && string(displayName[0]) == "." && string(displayName[1]) == "/" {
			displayName = displayName[2:]
		}
		renderedTabs = append(renderedTabs, style.Render(displayName))
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	_, h, _ := term.GetSize(0)
	doc.WriteString(windowStyle.Width((0)).Height(h - 100).Render(m.tabs[m.activeTab].View()))

	return docStyle.Render(doc.String())
}

func (m *Controller) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.tabs[m.activeTab].fileSelector.quitting {
		m.tabs[m.activeTab].fileSelector.Update(msg)
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "left":
			m.activeTab = max(0, m.activeTab-1)
		case "right":
			m.activeTab = min(len(m.files)-1, m.activeTab+1)
		default:
			m.tabs[m.activeTab].Update(msg)
		}
	}
	return m, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
