package cmd

import (
	"fmt"

	// Github imports
	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
)

func checkErr (err error) {
	if err != nil {
		panic(err)
	}
}

func createShit() {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "", // disables re-printing entered value
	}

	// Title input
	Ideaprompt := promptui.Prompt{
		Label:     "Enter your idea title",
		Templates: templates,
	}
	_, err := Ideaprompt.Run()
	checkErr(err)

	// Description input
	DescPrompt := promptui.Prompt{
		Label:     "Description for your idea son",
		Templates: templates,
	}
	_, descerr := DescPrompt.Run()
	checkErr(descerr)

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


