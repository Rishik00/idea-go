package teaui

import (
	"fmt"
	// "log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Copy pasted, sorry no frontend for me
var (
    focusedStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#afafafff")). // A soft indigo
        Padding(1, 2).
        MarginBottom(1).
        Bold(true).
        Background(lipgloss.Color("#060606ff")) // Dark background

    blurredStyle = lipgloss.NewStyle().
        Border(lipgloss.HiddenBorder()).
        Foreground(lipgloss.Color("#777")).
        Padding(1, 2).
        MarginBottom(1).
        Background(lipgloss.Color("#1A1A1A")) // Same as focused, just less intense

    helpTextStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#555")).
        Italic(true).
        MarginTop(1).
        Align(lipgloss.Center).
        Faint(true)

	mainTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#D7BA7D")). // warm gold-like color
		Align(lipgloss.Center).
		MarginBottom(2).
		Background(lipgloss.Color("#0F0F0F")).
		Padding(1, 0)
)


func initialAltInitModel() Initmodel {
    ti := textarea.New()
    ti.Placeholder = "Title"
    
	ti.SetHeight(2)
    ti.SetWidth(50)
    
	ti.Focus() // Title is focused first

    desc := textarea.New()
    desc.Placeholder = "Description"
    desc.SetHeight(10)
    desc.SetWidth(50)

    return Initmodel{
        title:       ti,
        description: desc,
        focusIndex:  0,
    }
}

func (m Initmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
			
        case "ctrl+c", "q":
            return m, tea.Quit

        case "tab", "shift+tab":
            m.focusIndex = (m.focusIndex + 1) % 2
            if m.focusIndex == 0 {
                m.title.Focus()
                m.description.Blur()
            } else {
                m.title.Blur()
                m.description.Focus()
            }
            return m, nil

        case "ctrl+s":
            fmt.Println("Saving: ", m.title.Value(), m.description.Value())
            return m, tea.Quit
        }
    }

    // This part was incorrectly placed — it should be here
    if m.focusIndex == 0 {
        m.title, cmd = m.title.Update(msg)
    } else {
        m.description, cmd = m.description.Update(msg)
    }

    return m, cmd
}


func (m Initmodel) View() string {
    var titleView, descView string

    // Apply focus styles
    if m.focusIndex == 0 {
        titleView = focusedStyle.Render(m.title.View())
        descView = blurredStyle.Render(m.description.View())
    } else {
        titleView = blurredStyle.Render(m.title.View())
        descView = focusedStyle.Render(m.description.View())
    }

    // Label style
    labelStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#9A9A9A")).
        MarginBottom(1)

    // Outer padding/margin around the full layout
    layoutPadding := lipgloss.NewStyle().
        MarginTop(2).
        MarginLeft(3).
        MarginRight(3)

    // Compose full layout
    layout := lipgloss.JoinVertical(lipgloss.Left,
		mainTitleStyle.Render(`
		╔════════════════════════════════════╗
		║       Description of title         ║
		╚════════════════════════════════════╝
		`),
        labelStyle.Render("Title"),
        titleView,
        labelStyle.Render("Description"),
        descView,
        helpTextStyle.Render("TAB to switch • CTRL+S to save • Q to quit"),
    )


    return layoutPadding.Render(layout)
}


func UseAltInit() (string, string, error) {
    p := tea.NewProgram(initialAltInitModel(), tea.WithAltScreen())
	finalAltModel, err := p.Run()

	if err != nil {
		panic(err)
	}

    if m, ok := finalAltModel.(Initmodel); ok {
        return m.title.Value(), m.description.Value(), err
    }

	return "No Title", "No Desc", err

}

func (m Initmodel) Init() tea.Cmd {
	return nil
}
