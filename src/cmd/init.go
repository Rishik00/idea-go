package cmd

import (
	"fmt"

	// Github imports
	"github.com/spf13/cobra"

	// Local imports
	"idea/db"
	"idea/teaui"
)
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func createBucket() []string {
	bucketName, err := teaui.UseTitle()
	checkErr(err)

	db.AddBucket(bucketName)
	setupPageId(bucketName)

	keys, _ := db.ShowExistingBuckets()
	return keys
}

func createShit() {
	keys, _ := db.ShowExistingBuckets()

	if len(keys) == 0 {
		fmt.Println("No workspaces found, so lets add one")
		keys = createBucket()
	}

	bucket, err := teaui.UseChoice(keys)
	checkErr(err)

	title, description, err := teaui.UseAltInit()
	checkErr(err)

	if title == "" || description == "" {
		fmt.Println("Either title or description or both are empty, please add them or else...")
		return;
	}

	fmt.Println("\nYour skill issues Idea has been added bub")
	db.AddIdea(bucket, title, description)
}

func InitBucketsAndIdeas() {
	var options = []string{"Workspace", "Title"}

	selected, err  := teaui.UseChoice(options)
	checkErr(err)
	
	if selected == "Add Workspace" {
		createBucket()
	} else {
		createShit()
	}
}

var InitCmd = &cobra.Command{
	Use: "init",
	Aliases: []string{"init"},
	Short: "Adding something",
	Run: func (cmd *cobra.Command, args []string){
		InitBucketsAndIdeas()
	},
}

func init() {
	rootCmd.AddCommand(InitCmd)
}