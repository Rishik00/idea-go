// The setup command will be done as soon as the package is done building up.
// Users can run: idea setup to see wassup

package cmd

import (
	"fmt"
	"os"

	// Github imports
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/joho/godotenv"

	// Local imports
	// "idea/db"
	// "idea/teaui"
)

func AddToEnv(k string , v string) (string, error) {
	// Load existing .env to preserve other keys
	existingEnv, _ := godotenv.Read(".env")
	existingEnv[k] = v

	err := godotenv.Write(existingEnv, ".env")
	if err != nil {
		return "Failure", err
	}

	return "Success", nil
}

// TODO: move this to bubbletea
func setupAPIKey() {
	api_key := ""
	Keyprompt := &survey.Input{
		Message: "Paste your API key:",
	}

	survey.AskOne(Keyprompt, &api_key)	
	_, err := AddToEnv("NOTION_API_KEY", api_key)

	checkErr(err)
}

func getRelevantKeys() ([]string) {
	existingEnv, _ := godotenv.Read(".env")

	var keys []string
	for k, v := range existingEnv {
		if v == "DEFAULTPARAM" {
			keys = append(keys, k)
		}
	}

	return keys
}

// TODO: move this to bubbletea
func setupPageId(pageName string) (string) {
	if pageName != "DEFAULTARGUMENT" {
		_, err := AddToEnv(pageName, "DEFAULTPARAM")
		checkErr(err)
		return "success"
	}
	
	rkeys := getRelevantKeys()

	var res string
	rKeyPrompt := &survey.Select{
		Message: "Choose a bucket:",
		Options: rkeys,
	}
	checkErr(survey.AskOne(rKeyPrompt, &res))

	page_id := ""
	PageIdPrompt := &survey.Input{
		Message: "Paste your page ID:",
	}
	survey.AskOne(PageIdPrompt, &page_id)

	_, err := AddToEnv(res, page_id)
	checkErr(err)

	return "Success"
}

func Setup() {
	fmt.Println("Checking your env variables for API keys and page IDs")
	err := godotenv.Load()
	checkErr(err)

	if os.Getenv("NOTION_API_KEY") == "" {
		fmt.Println("API KEY not found")
		setupAPIKey()
	} else {
		setupPageId("DEFAULTARGUMENT")
	}
}

var SetupCmd = &cobra.Command{
	Use: "setup",
	Aliases: []string{"setup"},
	Short: "Lets add some notion page keys and buckets. You only have to do this once or everytime you want to connect to a bucket.",
	Run: func(cmd *cobra.Command, args []string){
		Setup()
	},
}

func init() {
	rootCmd.AddCommand(SetupCmd)
}