package cmd

import (
	"fmt"

	// Github imports
	// "github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func Sync() {

}

var SyncCmd = &cobra.Command{
	Use: "sync",
	Aliases: []string{"sync"},
	Short: "Sync everything with my notion page with the notion API",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Print("Syncing everything with the cloud, ig")
		Sync()
	},
}


func init(){ 
	rootCmd.AddCommand(SyncCmd)
}