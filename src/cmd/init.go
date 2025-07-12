package cmd

import (
	"fmt"

	"idea/db"

	// Github imports
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func createBucket() []string {
	bucket_name := ""
	BucketPrompt := &survey.Input{
		Message: "Paste your bucket name:",
	}
	survey.AskOne(BucketPrompt, &bucket_name)

	db.AddBucket(bucket_name)

	// idk how this will go
	setupPageId(bucket_name)

	keys, _ := db.ShowExistingBuckets()
	return keys
}

func createShit() {
	keys, _ := db.ShowExistingBuckets()

	if len(keys) == 0 {
		fmt.Println("No buckets found, so lets add one")
		keys = createBucket()
	}

	fmt.Println(len(keys))

	var bucket string
	bucketPrompt := &survey.Select{
		Message: "Choose a bucket:",
		Options: keys,
	}
	checkErr(survey.AskOne(bucketPrompt, &bucket))

	var title string
	titlePrompt := &survey.Input{
		Message: "Enter your idea title:",
	}
	checkErr(survey.AskOne(titlePrompt, &title))

	var description string
	descPrompt := &survey.Input{
		Message: "Enter description for your idea:",
	}
	checkErr(survey.AskOne(descPrompt, &description))

	fmt.Println("\nYour Idea:")
	fmt.Println("Bucket     :", bucket)
	fmt.Println("Title      :", title)
	fmt.Println("Description:", description)

	db.AddIdea(bucket, title, description)
}

func InitBucketsAndIdeas() {
	var options = []string{"add bucket", "add idea"}
	
	var option string
	bucketPrompt := &survey.Select{
		Message: "Choose:",
		Options: options,
	}
	checkErr(survey.AskOne(bucketPrompt, &option))

	if option == "add bucket" {
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