package cmd

import (
	"log"

	"idea/tui"
)

func Execute() {
	if err := tui.Run(); err != nil {
		log.Fatalf("run tui: %v", err)
	}
}
