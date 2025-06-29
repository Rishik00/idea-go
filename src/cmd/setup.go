package cmd

import (
	"fmt"

	// Github imports
	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"

	// Local imports
	"idea/db"
)


func SetupShit() {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "", // disables re-printing entered value
	}

	// Title input
	ApiSecretPrompt := promptui.Prompt{
		Label:     "Enter your notion secret: ",
		Templates: templates,
	}
	_, err := ApiSecretPrompt.Run()
	checkErr(err)

}

func setupBucket() {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "", // disables re-printing entered value
	}

	// Title input
	ApiSecretPrompt := promptui.Prompt{
		Label:     "Enter your bucketname: ",
		Templates: templates,
	}
	result, err := ApiSecretPrompt.Run()
	checkErr(err)

	db.AddBucket(result)
	fmt.Print("Added the bucketname to DB:", result)
}

var SetupCmd = &cobra.Command{
	Use: "setup",
	Aliases: []string{"setup"},
	Short: "Just a few things for notion integration, you can do this once and forget about it.",
	Run: func (cmd *cobra.Command, args []string) {
		fmt.Print("Just a few things for notion integration, you can do this once and forget about it.")
		setupBucket()
	},
}

func init() {
	rootCmd.AddCommand(SetupCmd)
}
