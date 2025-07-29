package teaui

// A simple program demonstrating the textarea component from the Bubbles
// component library.

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func UseDescription() string {
	p := tea.NewProgram(initialDescriptionModel())

	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	if m, ok := finalModel.(Descriptionmodel); ok {
		return m.textarea.Value() 
	}

	return ""  // fallback in case of unexpected failure
}

func initialDescriptionModel() Descriptionmodel {
	ti := textarea.New()
	ti.Placeholder = "I want to bomb people, jk"
	ti.Focus()

	
	// Set size
	ti.SetWidth(200)  // any int
	ti.SetHeight(100) // number of lines visible

	return Descriptionmodel{
		textarea: ti,
		err:      nil,
	}
}

func (m Descriptionmodel) Init() tea.Cmd {
	return textarea.Blink
}

func (m Descriptionmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Descriptionmodel) View() string {
	return fmt.Sprintf(
		"Tell me a story.\n\n%s\n\n%s",
		m.textarea.View(),
		"(ctrl+c to quit)",
	) + "\n\n"
}