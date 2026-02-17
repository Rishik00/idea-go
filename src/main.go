package main

import (
	"log"

	"idea/cmd"
	"idea/db"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("init db: %v", err)
	}
	defer db.CloseDB()

	cmd.Execute()
}
