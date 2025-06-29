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

func createShit() {
	keys, err := db.ShowExistingBuckets()
	checkErr(err)

	if len(keys) == 0 {
		fmt.Println("No buckets found.")
		return
	}

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

}

var InitCmd = &cobra.Command{
	Use: "init",
	Aliases: []string{"init"},
	Short: "Adding an idea and a description for your ease",
	Run: func (cmd *cobra.Command, args []string){
		fmt.Print("Init command run hahaha")
		createShit()
	},
}

func init() {
	rootCmd.AddCommand(InitCmd)
}