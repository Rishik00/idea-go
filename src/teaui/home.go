package teaui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ASCII_ART = `
  ___ ____  _____   _    ____ _     ___ 
 |_ _|  _ \| ____| / \  / ___| |   |_ _|
  | || | | |  _|  / _ \| |   | |    | | 
  | || |_| | |___/ ___ \ |___| |___ | | 
 |___|____/|_____/_/   \_\____|_____|___|

 Why are you here? Anyway hi :)
 ----------------------------
 Use 'idea init' to create a new idea
 Use 'idea --help' to see all commands
`

// Style
var (
	bgStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#0a0a0a")).
		Foreground(lipgloss.Color("#cccccc")).
		Padding(2, 4)

	asciiStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ffcc")).
		Italic(true).
		Bold(true).
		MarginBottom(1)
)

// Model
type homeModel struct{}

func (m homeModel) Init() tea.Cmd {
	return nil
}

func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m homeModel) View() string {
	ascii := asciiStyle.Render(ASCII_ART)
	return bgStyle.Render(ascii + "\n\nPress Q or ESC to quit.")
}

// Public function to start home view
func UseHomePage() {
	p := tea.NewProgram(homeModel{})
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
