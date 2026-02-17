package syncer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"idea/db"

	"github.com/joho/godotenv"
)

const defaultObsidianPath = "~/.idea/ideas.md"

func SyncFromDB() error {
	ideas, err := db.ListIdeas()
	if err != nil {
		return err
	}
	return WriteIdeasMarkdown(ResolveObsidianPath(), ideas)
}

func ResolveObsidianPath() string {
	_ = godotenv.Load()

	path := strings.TrimSpace(os.Getenv("OBSIDIAN_PATH"))
	if path == "" {
		path = defaultObsidianPath
	}

	return expandHome(path)
}

func WriteIdeasMarkdown(path string, ideas []db.StoredIdea) error {
	sorted := make([]db.StoredIdea, len(ideas))
	copy(sorted, ideas)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Idea.CreatedAt == sorted[j].Idea.CreatedAt {
			return sorted[i].ID > sorted[j].ID
		}
		return sorted[i].Idea.CreatedAt > sorted[j].Idea.CreatedAt
	})

	var b strings.Builder
	b.WriteString("# Ideas\n\n")

	for _, item := range sorted {
		title := strings.TrimSpace(item.Idea.Title)
		if title == "" {
			title = "Untitled idea"
		}

		ts := time.Unix(item.Idea.CreatedAt, 0).Format("2006-01-02 15:04")
		raw := markdownValue(item.Idea.Raw, "(empty)")
		elaborated := markdownValue(item.Idea.Elaborated, "(pending elaboration)")

		fmt.Fprintf(&b, "## %s\n", title)
		fmt.Fprintf(&b, "*%s*\n\n", ts)
		fmt.Fprintf(&b, "**Raw:** %s\n\n", raw)
		fmt.Fprintf(&b, "**Elaborated:** %s\n\n", elaborated)
		b.WriteString("---\n\n")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0644)
}

func markdownValue(v string, fallback string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(v, "\r\n", "\n"))
	if clean == "" {
		return fallback
	}
	return clean
}

func expandHome(path string) string {
	if path == "~" {
		return userHomeDir()
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		return filepath.Join(userHomeDir(), path[2:])
	}
	return path
}

func userHomeDir() string {
	if h, ok := os.LookupEnv("USERPROFILE"); ok && h != "" {
		return h
	}
	if h, ok := os.LookupEnv("HOME"); ok && h != "" {
		return h
	}
	return "."
}
