package tui

import (
	"fmt"
	"strings"
	"time"
)

func renderIdeaView(m appModel) string {
	width := m.contentWidth()

	item, ok := m.viewingIdea()
	if !ok {
		return bodyStyle.Width(width).Render(dangerStyle.Render("Idea not found. Press esc to return."))
	}

	title := strings.TrimSpace(item.Idea.Title)
	if title == "" {
		title = "Untitled idea"
	}
	raw := strings.TrimSpace(item.Idea.Raw)
	if raw == "" {
		raw = "(empty)"
	}
	elaborated := strings.TrimSpace(item.Idea.Elaborated)
	elabStatus := "[elaborating...]"
	if elaborated == "" {
		elaborated = "(pending elaboration)"
	} else {
		elabStatus = "[ready]"
	}

	created := time.Unix(item.Idea.CreatedAt, 0).Format("2006-01-02 15:04")

	lines := []string{
		titleLabelStyle.Render(title),
		separatorStyle.Render(strings.Repeat("-", width-2)),
		mutedStyle.Render(fmt.Sprintf("Created: %s", created)),
		"",
		titleLabelStyle.Render("Raw"),
		raw,
		"",
		titleLabelStyle.Render("Elaborated (LLM)"),
		elaborated,
		"",
		mutedStyle.Render(elabStatus),
	}

	if m.confirmDelete {
		lines = append(lines, "")
		lines = append(lines, dangerStyle.Render("Delete this idea? (y/n)"))
	}

	return bodyStyle.Width(width).Render(strings.Join(lines, "\n"))
}
