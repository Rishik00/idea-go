/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	// Local imports
	"idea/db"
)

var ASCII_ART = `
  ___ ____  _____    _    
 |_ _|  _ \| ____|  / \   
  | || | | |  _|   / _ \  
  | || |_| | |___ / ___ \ 
 |___|____/|_____/_/   \_\

 Why are you here? Anyway hi :)
 ----------------------------
 Use 'idea init' to create a new idea
 Use 'idea --help' to see all commands
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "idea",
	Short: "Just an idea management CLI",
	Long:  `Just another idea management CLI with git-like commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ASCII_ART)
		db.PrintExistingBuckets()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.idea.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


