package synclayer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

// Yes, I took GPT's help here. Maybe I should write some tutorials
// About JSON  marshalling and unmarshalling.
type NotionPage struct {
	Parent     Parent     `json:"parent"`
	Properties Properties `json:"properties"`
	Children   []Block    `json:"children"`
}

type Parent struct {
	Type   string `json:"type"`
	PageID string `json:"page_id"`
}

type Properties struct {
	Title Title `json:"title"`
}

type Title struct {
	Title []TextObject `json:"title"`
}


type Block struct {
	Object    string    `json:"object"`
	Type      string    `json:"type"`
	Paragraph Paragraph `json:"paragraph"`
}

type Paragraph struct {
	RichText []TextObject `json:"rich_text"`
}

type TextObject struct {
	Type string  `json:"type"`
	Text Text    `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}


func GetTokens() (string, string) {
	err := godotenv.Load()
	checkErr(err)

	fmt.Print("Page ID: ", os.Getenv("PAGE_ID_ONE"), "\n")
	fmt.Print("Key: ", os.Getenv("NOTION_API_KEY"), "\n")

	return os.Getenv("PAGE_ID_ONE"),os.Getenv("NOTION_API_KEY")
}

func FormatBody(pageid, title, description string) ([]byte, error) {
	children := []Block{
		{
			Object: "block",
			Type:   "paragraph",
			Paragraph: Paragraph{
				RichText: []TextObject{
					{
						Type: "text",
						Text: Text{
							Content: description,
						},
					},
				},
			},
		},
	}

	parent := Parent{
		Type:   "page_id",
		PageID: pageid,
	}

	properties := Properties{
		Title: Title{
			Title: []TextObject{
				{
					Type: "text",
					Text: Text{
						Content: title,
					},
				},
			},
		},
	}

	JSONStuff := NotionPage{
		Parent: parent,
		Properties: properties,
		Children: children,
	}

	return json.Marshal(JSONStuff)
}


func PostIdea(title, description string) {
	PageId, APiKey := GetTokens()
	fmt.Print(PageId)

	NotionPayload, err := FormatBody(PageId, title, description)
	checkErr(err)

	client := &http.Client{}
	Endpoint := "https://api.notion.com/v1/pages/"
	
	req, err := http.NewRequest("POST", Endpoint, bytes.NewBuffer(NotionPayload))
	checkErr(err)

	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization",  "Bearer "+ APiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	checkErr(err)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println("❌ Error response from Notion:", string(bodyBytes))
	}
	fmt.Print("Stuff sent with status code: ", resp.StatusCode, "\n")
}