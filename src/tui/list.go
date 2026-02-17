package tui

import (
	"fmt"
	"strings"
	"time"

	"idea/db"
)

func renderList(m appModel) string {
	ideas := m.currentIdeas()
	width := m.contentWidth()
	lines := make([]string, 0, 24)

	if m.mode == modeSearch {
		lines = append(lines, titleLabelStyle.Render("Search"))
		lines = append(lines, inputFocusedStyle.Width(width-2).Render(m.searchInput.View()))
		lines = append(lines, "")
	}

	if len(ideas) == 0 {
		if m.mode == modeSearch {
			lines = append(lines, mutedStyle.Render("No matching ideas. Press esc to go back."))
		} else {
			lines = append(lines, mutedStyle.Render("No ideas yet. Press n to add your first one."))
		}
		return bodyStyle.Width(width).Render(strings.Join(lines, "\n"))
	}

	start, end := m.visibleWindow(len(ideas))
	for i := start; i < end; i++ {
		marker := "  "
		style := normalItemStyle
		if i == m.cursor {
			marker = "> "
			style = selectedItemStyle
		}

		title := strings.TrimSpace(ideas[i].Idea.Title)
		if title == "" {
			title = "Untitled idea"
		}
		title = truncate(title, width-14)

		ts := time.Unix(ideas[i].Idea.CreatedAt, 0).Format("01-02 15:04")
		line := fmt.Sprintf("%s%s %s", marker, title, mutedStyle.Render(ts))
		lines = append(lines, style.Render(line))
	}

	if m.confirmDelete {
		lines = append(lines, "")
		lines = append(lines, dangerStyle.Render("Delete this idea? (y/n)"))
	}

	return bodyStyle.Width(width).Render(strings.Join(lines, "\n"))
}

func renderHeader(m appModel) string {
	title := "idea"
	switch m.mode {
	case modeInsert:
		if m.editingID != 0 {
			title = "idea - edit"
		} else {
			title = "idea - new"
		}
	case modeSearch:
		title = "idea - search"
	case modeView:
		title = "idea - view"
	}

	modeLabel := fmt.Sprintf("[%s]", strings.ToUpper(string(m.mode)))
	left := fmt.Sprintf(" %s", title)
	pad := m.contentWidth() - len(left) - len(modeLabel)
	if pad < 1 {
		pad = 1
	}

	return headerStyle.Width(m.contentWidth()).Render(left + strings.Repeat(" ", pad) + modeLabel)
}

func renderFooter(m appModel) string {
	text := ""
	switch m.mode {
	case modeNormal:
		text = "n new   / search   d delete   v view   q quit"
	case modeSearch:
		text = "enter view   j/k navigate   esc back"
	case modeInsert:
		text = "tab switch field   ctrl+s save   esc cancel   space voice (phase 3)"
	case modeView:
		text = "e edit   d delete   esc back"
	}

	return footerStyle.Width(m.contentWidth()).Render(text)
}

func (m appModel) visibleWindow(total int) (int, int) {
	maxRows := 12
	if m.height > 0 {
		maxRows = m.height - 12
	}
	if maxRows < 5 {
		maxRows = 5
	}

	if total <= maxRows {
		return 0, total
	}

	start := m.cursor - (maxRows / 2)
	if start < 0 {
		start = 0
	}
	end := start + maxRows
	if end > total {
		end = total
		start = end - maxRows
	}
	return start, end
}

func filterIdeas(ideas []db.StoredIdea, query string) []db.StoredIdea {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return ideas
	}

	filtered := make([]db.StoredIdea, 0, len(ideas))
	for _, idea := range ideas {
		title := strings.ToLower(idea.Idea.Title)
		raw := strings.ToLower(idea.Idea.Raw)
		if strings.Contains(title, q) || strings.Contains(raw, q) {
			filtered = append(filtered, idea)
		}
	}
	return filtered
}

func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}

	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n <= 3 {
		return string(r[:n])
	}
	return string(r[:n-3]) + "..."
}
