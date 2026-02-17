package tui

import "github.com/charmbracelet/lipgloss"

var (
	appFrameStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#6C6EA0")).
		Padding(0, 1).
		Background(lipgloss.Color("#0F0F0F"))

	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#D7BA7D")).
		Padding(0, 0, 1, 0)

	bodyStyle = lipgloss.NewStyle().
		Padding(0, 0, 1, 0).
		Foreground(lipgloss.Color("#E5E5E5"))

	footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA0B2")).
		BorderTop(true).
		BorderForeground(lipgloss.Color("#6C6EA0")).
		Padding(1, 0, 0, 0)

	selectedItemStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F4E0A3"))

	normalItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CFCFD9"))

	mutedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B8F9F"))

	titleLabelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D7BA7D")).
		Bold(true)

	inputFocusedStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#6C6EA0")).
		Padding(0, 1)

	inputBlurredStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3F4250")).
		Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A7D59B")).
		Padding(0, 0, 0, 0)

	dangerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F08C8C")).
		Bold(true)

	separatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6C6EA0"))
)
