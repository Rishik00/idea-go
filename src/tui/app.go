package tui

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"idea/db"
	syncer "idea/sync"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode string

const (
	modeNormal mode = "normal"
	modeInsert mode = "insert"
	modeSearch mode = "search"
	modeView   mode = "view"
)

type appModel struct {
	mode mode

	width  int
	height int

	ideas  []db.StoredIdea
	cursor int

	searchInput textinput.Model

	titleInput textinput.Model
	rawInput   textarea.Model

	composeFocus int
	editingID    uint64

	viewingID uint64

	confirmDelete bool
	status        string
	errorText     string
}

func Run() error {
	model, err := newModel()
	if err != nil {
		return err
	}

	program := tea.NewProgram(model, tea.WithAltScreen())
	_, err = program.Run()
	return err
}

func newModel() (appModel, error) {
	ideas, err := db.ListIdeas()
	if err != nil {
		return appModel{}, err
	}
	status := fmt.Sprintf("%d ideas loaded", len(ideas))
	if syncErr := syncer.SyncFromDB(); syncErr != nil {
		status = fmt.Sprintf("%d ideas loaded (obsidian sync failed)", len(ideas))
	}

	m := appModel{
		mode:        modeNormal,
		ideas:        ideas,
		searchInput:  newSearchInput(),
		titleInput:   newTitleInput(),
		rawInput:     newRawInput(),
		status:       status,
		errorText:    "",
		composeFocus: 0,
	}
	m.setComposeFocus(0)
	return m, nil
}

func newSearchInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "filter ideas..."
	ti.CharLimit = 240
	ti.Width = 80
	ti.Blur()
	return ti
}

func newTitleInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Idea title"
	ti.CharLimit = 200
	ti.Width = 80
	ti.Focus()
	return ti
}

func newRawInput() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Write the raw thought..."
	ta.SetHeight(8)
	ta.SetWidth(80)
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	ta.Blur()
	return ta
}

func (m appModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.rawInput.SetWidth(m.contentWidth() - 6)
		m.titleInput.Width = m.contentWidth() - 6
		m.searchInput.Width = m.contentWidth() - 6
	}

	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	switch m.mode {
	case modeInsert:
		return m.updateInsert(msg)
	case modeSearch:
		return m.updateSearch(msg)
	case modeView:
		return m.updateView(msg)
	default:
		return m.updateNormal(msg)
	}
}

func (m appModel) updateNormal(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.confirmDelete {
		switch key.String() {
		case "y":
			selected, ok := m.selectedIdea()
			if !ok {
				m.confirmDelete = false
				return m, nil
			}
			m = m.deleteIdea(selected.ID)
			return m, nil
		case "n", "esc":
			m.confirmDelete = false
			m.status = "delete cancelled"
			return m, nil
		}
	}

	switch key.String() {
	case "q":
		return m, tea.Quit
	case "j", "down":
		m.cursor++
		m.clampCursor()
	case "k", "up":
		m.cursor--
		m.clampCursor()
	case "n":
		m = m.startInsert()
	case "/":
		m = m.startSearch()
	case "v", "enter":
		m = m.openSelectedView()
	case "d":
		if _, ok := m.selectedIdea(); ok {
			m.confirmDelete = true
			m.status = "confirm delete (y/n)"
		}
	}

	return m, nil
}

func (m appModel) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc":
			m.mode = modeNormal
			m.searchInput.Blur()
			m.searchInput.SetValue("")
			m.cursor = 0
			m.status = fmt.Sprintf("%d ideas loaded", len(m.ideas))
			m.clampCursor()
			return m, nil
		case "j", "down":
			m.cursor++
			m.clampCursor()
			return m, nil
		case "k", "up":
			m.cursor--
			m.clampCursor()
			return m, nil
		case "enter", "v":
			m = m.openSelectedView()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	m.cursor = 0
	m.clampCursor()
	count := len(m.currentIdeas())
	m.status = fmt.Sprintf("%d matches", count)
	return m, cmd
}

func (m appModel) updateInsert(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc":
			m.mode = modeNormal
			m.editingID = 0
			m.titleInput = newTitleInput()
			m.rawInput = newRawInput()
			m.setComposeFocus(0)
			m.status = "compose cancelled"
			return m, nil
		case "tab", "shift+tab":
			m.setComposeFocus(m.composeFocus + 1)
			return m, nil
		case "ctrl+s":
			m = m.saveCompose()
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.composeFocus == 0 {
		m.titleInput, cmd = m.titleInput.Update(msg)
	} else {
		m.rawInput, cmd = m.rawInput.Update(msg)
	}
	return m, cmd
}

func (m appModel) updateView(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.confirmDelete {
		switch key.String() {
		case "y":
			m = m.deleteIdea(m.viewingID)
			m.mode = modeNormal
			return m, nil
		case "n", "esc":
			m.confirmDelete = false
			m.status = "delete cancelled"
			return m, nil
		}
	}

	switch key.String() {
	case "esc":
		m.mode = modeNormal
		m.confirmDelete = false
	case "e":
		m = m.startEditFromView()
	case "d":
		if _, ok := m.viewingIdea(); ok {
			m.confirmDelete = true
			m.status = "confirm delete (y/n)"
		}
	}

	return m, nil
}

func (m appModel) View() string {
	var body string
	switch m.mode {
	case modeInsert:
		body = renderCompose(m)
	case modeView:
		body = renderIdeaView(m)
	default:
		body = renderList(m)
	}

	parts := []string{
		renderHeader(m),
		body,
		renderFooter(m),
	}
	if strings.TrimSpace(m.status) != "" {
		parts = append(parts, statusStyle.Width(m.contentWidth()).Render(m.status))
	}
	if strings.TrimSpace(m.errorText) != "" {
		parts = append(parts, dangerStyle.Width(m.contentWidth()).Render(m.errorText))
	}

	panel := appFrameStyle.Width(m.contentWidth() + 2).Render(lipgloss.JoinVertical(lipgloss.Left, parts...))
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, panel)
	}
	return panel
}

func (m appModel) startInsert() appModel {
	m.mode = modeInsert
	m.editingID = 0
	m.titleInput = newTitleInput()
	m.rawInput = newRawInput()
	m.setComposeFocus(0)
	m.confirmDelete = false
	m.errorText = ""
	m.status = "new idea"
	return m
}

func (m appModel) startSearch() appModel {
	m.mode = modeSearch
	m.searchInput = newSearchInput()
	m.searchInput.Focus()
	m.cursor = 0
	m.confirmDelete = false
	m.errorText = ""
	m.status = "search mode"
	return m
}

func (m appModel) openSelectedView() appModel {
	selected, ok := m.selectedIdea()
	if !ok {
		m.status = "no idea selected"
		return m
	}

	m.viewingID = selected.ID
	m.mode = modeView
	m.confirmDelete = false
	m.errorText = ""
	return m
}

func (m appModel) startEditFromView() appModel {
	item, ok := m.viewingIdea()
	if !ok {
		m.mode = modeNormal
		m.errorText = "idea no longer exists"
		return m
	}

	m.mode = modeInsert
	m.editingID = item.ID
	m.titleInput = newTitleInput()
	m.rawInput = newRawInput()
	m.titleInput.SetValue(item.Idea.Title)
	m.rawInput.SetValue(item.Idea.Raw)
	m.setComposeFocus(0)
	m.status = "editing idea"
	m.errorText = ""
	return m
}

func (m appModel) saveCompose() appModel {
	title := strings.TrimSpace(m.titleInput.Value())
	raw := strings.TrimSpace(m.rawInput.Value())
	if title == "" {
		m.status = ""
		m.errorText = "title is required"
		return m
	}
	if raw == "" {
		raw = title
	}

	if m.editingID != 0 {
		current, err := db.GetIdea(m.editingID)
		if err != nil {
			m.errorText = err.Error()
			return m
		}

		next := current.Idea
		next.Title = title
		next.Raw = raw

		if err := db.UpdateIdea(m.editingID, next); err != nil {
			m.errorText = err.Error()
			return m
		}
		m.status = "idea updated"
	} else {
		if _, err := db.PutIdea(db.Idea{
			Title:      title,
			Raw:        raw,
			Elaborated: "",
			CreatedAt:  time.Now().Unix(),
		}); err != nil {
			m.errorText = err.Error()
			return m
		}
		m.status = "idea saved"
	}

	if err := syncer.SyncFromDB(); err != nil {
		m.errorText = fmt.Sprintf("saved, but obsidian sync failed: %v", err)
	} else {
		m.errorText = ""
	}

	m = m.reloadIdeas()
	m.mode = modeNormal
	m.editingID = 0
	m.titleInput = newTitleInput()
	m.rawInput = newRawInput()
	m.setComposeFocus(0)
	return m
}

func (m appModel) deleteIdea(id uint64) appModel {
	if err := db.DeleteIdea(id); err != nil {
		if errors.Is(err, db.ErrIdeaNotFound) {
			m.errorText = "idea not found"
		} else {
			m.errorText = err.Error()
		}
		m.confirmDelete = false
		return m
	}

	if err := syncer.SyncFromDB(); err != nil {
		m.errorText = fmt.Sprintf("deleted, but obsidian sync failed: %v", err)
	} else {
		m.errorText = ""
	}

	m = m.reloadIdeas()
	m.confirmDelete = false
	m.status = "idea deleted"
	return m
}

func (m appModel) reloadIdeas() appModel {
	ideas, err := db.ListIdeas()
	if err != nil {
		m.errorText = err.Error()
		return m
	}

	m.ideas = ideas
	m.clampCursor()
	return m
}

func (m appModel) selectedIdea() (db.StoredIdea, bool) {
	ideas := m.currentIdeas()
	if len(ideas) == 0 {
		return db.StoredIdea{}, false
	}
	if m.cursor < 0 || m.cursor >= len(ideas) {
		return db.StoredIdea{}, false
	}
	return ideas[m.cursor], true
}

func (m appModel) viewingIdea() (db.StoredIdea, bool) {
	for _, item := range m.ideas {
		if item.ID == m.viewingID {
			return item, true
		}
	}
	return db.StoredIdea{}, false
}

func (m appModel) currentIdeas() []db.StoredIdea {
	if m.mode == modeSearch {
		return filterIdeas(m.ideas, m.searchInput.Value())
	}
	return m.ideas
}

func (m *appModel) clampCursor() {
	ideas := m.currentIdeas()
	if len(ideas) == 0 {
		m.cursor = 0
		return
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(ideas) {
		m.cursor = len(ideas) - 1
	}
}

func (m appModel) contentWidth() int {
	if m.width <= 0 {
		return 88
	}

	width := m.width - 8
	if width < 64 {
		return 64
	}
	if width > 120 {
		return 120
	}
	return width
}
