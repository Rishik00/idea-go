package tui

import (
	"strings"
)

func renderCompose(m appModel) string {
	width := m.contentWidth()

	titleBoxStyle := inputBlurredStyle
	rawBoxStyle := inputBlurredStyle
	if m.composeFocus == 0 {
		titleBoxStyle = inputFocusedStyle
	} else {
		rawBoxStyle = inputFocusedStyle
	}

	lines := []string{
		titleLabelStyle.Render("Title"),
		titleBoxStyle.Width(width - 2).Render(m.titleInput.View()),
		"",
		titleLabelStyle.Render("Raw"),
		rawBoxStyle.Width(width - 2).Render(m.rawInput.View()),
		"",
		mutedStyle.Render("* REC (space to stop, esc to cancel)"),
	}

	if strings.TrimSpace(m.status) != "" {
		lines = append(lines, "")
		lines = append(lines, statusStyle.Render(m.status))
	}

	return bodyStyle.Width(width).Render(strings.Join(lines, "\n"))
}

func (m *appModel) setComposeFocus(next int) {
	if next < 0 {
		next = 1
	}
	if next > 1 {
		next = 0
	}
	m.composeFocus = next

	if m.composeFocus == 0 {
		m.titleInput.Focus()
		m.rawInput.Blur()
		return
	}

	m.titleInput.Blur()
	m.rawInput.Focus()
}
