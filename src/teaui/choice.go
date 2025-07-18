package teaui

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func UseChoice(options []string) (string, error) {
	p := tea.NewProgram(InitialChoiceModel(options))

	// Run returns the model as a tea.Model.
	m, err := p.Run()
	if err != nil {
		return "Problem in the choice space", err
	}
	mo := m.(Choicemodel)

	return mo.Choice, nil
}


func InitialChoiceModel(choices []string) Choicemodel {
    return Choicemodel{
        Choices:  choices,
    }
}

func (m Choicemodel) Init() tea.Cmd {
	return nil
}

func (m Choicemodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.Choice = m.Choices[m.Cursor]
			return m, tea.Quit

		case "down", "j":
			m.Cursor++
			if m.Cursor >= len(m.Choices) {
				m.Cursor = 0
			}

		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Choices) - 1
			}
		}
	}

	return m, nil
}

func (m Choicemodel) View() string {
	s := strings.Builder{}

	// Define colors
	const (
		cursorStyle  = "\033[1;35m" // Bright purple
		choiceStyle  = "\033[1;37m" // Bright gray/white
		resetStyle   = "\033[0m"
	)

	s.WriteString("\033[1;37mHello bub, what do you wanna do?\033[0m\n\n\n")

	for i := 0; i < len(m.Choices); i++ {
		if m.Cursor == i {
			s.WriteString(cursorStyle)
			s.WriteString("➤ ") // Arrow as cursor
			s.WriteString(m.Choices[i])
			s.WriteString(resetStyle)
			s.WriteString("\n")
		} else {
			s.WriteString("  ") // Padding for alignment
			s.WriteString(choiceStyle)
			s.WriteString(m.Choices[i])
			s.WriteString(resetStyle)
			s.WriteString("\n")
		}
	}

	s.WriteString("\n\n\033[2m(Use ↑/↓ and press Enter to select. q to quit.)\033[0m\n") // Dim style

	return s.String()
}
