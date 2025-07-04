package cmd

import (
	"fmt"

	// Github imports	
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	// Local imports
	"idea/synclayer"
	"idea/db"
)


func Sync() {
	var choices = []string{"All", "Select"}

	var choice string 
	choicePrompt := &survey.Select{
		Message: 	"Choose: ",
		Options: 	choices,
	}
	checkErr(survey.AskOne(choicePrompt, &choice))
	
	fmt.Print("Chosen choice: ", choice, "\n")
	if choice == "All" {
		synclayer.PostIdea("new Idea", "New New Idea")
	} else {
		keys, err := db.ShowExistingBuckets()
		checkErr(err)

		var bucket string
		bucketPrompt := &survey.Select{
			Message: "Choose a bucket:",
			Options: keys,
		}
		checkErr(survey.AskOne(bucketPrompt, &bucket))

		fmt.Print("Bucket selected: ", bucket, "\n")
		bucketIdeas, err := db.IdeasPerBucket(bucket)
		checkErr(err)

		for i := 0; i<len(bucketIdeas.Title); i+=1 {
			synclayer.PostIdea(bucketIdeas.Title[i], bucketIdeas.Desc[i])
		}
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