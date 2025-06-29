/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"idea/cmd"
	"idea/db"
)

func main() {
	err := db.InitDB() // opens my.db once
	if err != nil {
		panic(err)
	}
	defer db.CloseDB() // closes when app exits

	cmd.Execute()
}
