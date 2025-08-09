package cmd

import (
	// "fmt"
	// "os"

	// Github imports
	"github.com/spf13/cobra"
	// "github.com/joho/godotenv"

	// Local imports
	"idea/synclayer"
	// "idea/teaui"
)

func Pull() {
	synclayer.PullFromNotion()
}

var pullCmd = &cobra.Command{
	Use: "pull",
	Aliases: []string{"pull"},
	Short: "Pull from existing notion page of ideas",
	Run: func(cmd *cobra.Command, args []string){
		Pull()
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}