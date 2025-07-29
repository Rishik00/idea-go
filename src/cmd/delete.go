package cmd

import (
	// "fmt"

	// Github imports
	"fmt"

	"github.com/spf13/cobra"

	// Local imports
	// "idea/db"
	"idea/db"
	"idea/teaui"
)


func deleteInit() {
	var options = []string{"Workspace", "Title"}

	selected, err  := teaui.UseChoice(options)
	checkErr(err)

	if selected == "Workspace" {
		keys, _ := db.ShowExistingBuckets()

		if len(keys) == 0 {
			fmt.Println("Lol, what're you deleting when you have nothing?")
			keys = createBucket()
		}

		bucket, err := teaui.UseChoice(keys)
		checkErr(err)

		db.DeleteBucket(bucket)
	} else {
		fuzzySearchTitle(selected)
	}
}

func fuzzySearchTitle(title string) (string, error) {
	return "found", nil
}

var DeleteCmd = &cobra.Command{
	Use: "delete",
	Aliases: []string{"delete"},
	Short: "Deleting idea/workspace",
	Run: func (cmd *cobra.Command, args []string){
		deleteInit()		
	},
}

func init() {
	rootCmd.AddCommand(DeleteCmd)
}