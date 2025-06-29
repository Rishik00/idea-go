package cmd

import (
	"fmt"

	// Github imports	
	// survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func Pull () {
	
}

var pullCmd = &cobra.Command{
	Use: "pull",
	Aliases: []string{"pull"},
	Short: "Pull N number of your stored ideas and edit them?",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Print("Pull N number of your stored ideas and edit them?")
	},
}


func init(){ 
	rootCmd.AddCommand(pullCmd)
}