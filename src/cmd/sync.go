package cmd

import (
	"fmt"

	// Github imports	
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	// Local imports
	"idea/synclayer"
)


func Sync() {
	var choices = []string{"All of the available ones", "Select few"}

	var choice string 
	choicePrompt := &survey.Select{
		Message: 	"Choose: ",
		Options: 	choices,
	}
	checkErr(survey.AskOne(choicePrompt, &choice))
	
	fmt.Print("Chosen choice: ", choice, "\n")
	if choice == "All of the available ones" {
		synclayer.PostIdea("new Idea", "New New Idea")
	}
}

var SyncCmd = &cobra.Command{
	Use: "syncup",
	Aliases: []string{"syncup"},
	Short: "Sync everything with my notion page with the notion API",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Print("Syncing everything with the cloud, ig")
		Sync()
	},
}


func init(){ 
	rootCmd.AddCommand(SyncCmd)
}