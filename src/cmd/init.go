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
		fmt.Println("No buckets found, so lets add one")
		keys = createBucket()
	}

	bucket, err := teaui.UseChoice(keys)
	checkErr(err)

	title, _ := teaui.UseTitle()
	description := teaui.UseDescription()

	fmt.Println("\nYour Idea:")
	fmt.Println("Bucket     :", bucket)
	fmt.Println("Title      :", title)
	fmt.Println("Description:", description)

	db.AddIdea(bucket, title, description)
}

func InitBucketsAndIdeas() {
	var options = []string{"add bucket", "add idea"}

	selected, err  := teaui.UseChoice(options)
	checkErr(err)
	
	if selected == "add bucket" {
		createBucket()
	} else {
		createShit()
	}

}

var InitCmd = &cobra.Command{
	Use: "init",
	Aliases: []string{"init"},
	Short: "Adding an idea and a description for your ease",
	Run: func (cmd *cobra.Command, args []string){
		fmt.Print("Init command run hahaha")
		InitBucketsAndIdeas()
	},
}

func init() {
	rootCmd.AddCommand(InitCmd)
}