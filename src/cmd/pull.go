package cmd

import (
	"fmt"
	"idea/db"

	// Github imports
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	// Local imports
)

func PullLocal() error {
	keys, err := db.ShowExistingBuckets()
	checkErr(err)

	if len(keys) == 0 {
		fmt.Println("No buckets found.")
		return nil
	}

	var bucket string
	bucketPrompt := &survey.Select{
		Message: "Choose a bucket:",
		Options: keys,
	}
	checkErr(survey.AskOne(bucketPrompt, &bucket))

	fmt.Print("Bucket selected: ", bucket, "\n")
	bucketodIdeas, err := db.IdeasPerBucket(bucket)
	checkErr(err)

	var ideaAns string
	IdeaPrompt := &survey.Select{
		Message: "Choose a Idea: ",
		Options: bucketodIdeas.Title,
		// Description: func(title string, desc int) string {
		// 	var title_idx int
		// 	for idx, value := range bucketodIdeas.Title {
		// 		if value == title {
		// 			title_idx = idx
		// 			break
		// 		}
		// 	}
		// 	return bucketodIdeas.Desc[title_idx]
		// },
	}

	checkErr(survey.AskOne(IdeaPrompt, &ideaAns))
	
	return nil
}

func PullNotion() {

}

func PullEntryPt() {
	options := []string{"local", "notion"}
	
	var pullOption string
	pullPrompt := &survey.Select{
		Message: "Local or Notion: ",
		Options: options,
	}
	checkErr(survey.AskOne(pullPrompt, &pullOption))

	fmt.Print("Selected option: ", pullOption)
	if pullOption == "local" {
		PullLocal()
	}

}

var pullCmd = &cobra.Command{
	Use: "pull",
	Aliases: []string{"pull"},
	Short: "Pull N number of your stored ideas and edit them?",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Print("Pull N number of your stored ideas and edit them?")
		PullEntryPt()
	},
}


func init(){ 
	rootCmd.AddCommand(pullCmd)
}