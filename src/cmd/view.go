package cmd

import (
	"fmt"

	// Github imports
	"github.com/spf13/cobra"
	"github.com/charmbracelet/bubbles/list"
	
	// Local imports
	// "idea/synclayer"
	"idea/db"
	"idea/teaui"
)

type IdeaItem struct {
    Idea db.Idea	
}

func (i IdeaItem) Title() string        { return i.Idea.Title }
func (i IdeaItem) Description() string  { return i.Idea.Desc }
func (i IdeaItem) FilterValue() string  { return i.Idea.Title }


func viewFunc() error {
	keys, err := db.ShowExistingBuckets()
	checkErr(err)

	if len(keys) == 0 {
		fmt.Println("No buckets found.")
		return nil
	}
	
	bucket, err := teaui.UseChoice(keys)

	fmt.Print("Bucket selected: ", bucket, "\n")
	bucketofIdeas, err := db.IdeasPerBucket(bucket)
	checkErr(err)

	var items []list.Item
	for _, idea := range bucketofIdeas {
		items = append(items, IdeaItem{Idea: idea})
	}

	teaui.UseList(items, bucket)
	return nil
}

var viewCmd = &cobra.Command{
	Use: "view",
	Aliases: []string{"view"},
	Short: "Show your ideas and edit them?",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Print("Show your ideas and edit them?")
		viewFunc()
	},
}


func init(){ 
	rootCmd.AddCommand(viewCmd)
}