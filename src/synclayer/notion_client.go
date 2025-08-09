// MEGATODO: rewrite this, its dogpoop written by GPT 
// I hate JSON
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
		fmt.Println("Error response from Notion:", string(bodyBytes))
	}
	fmt.Print("Ideas sent with status code: ", resp.StatusCode, "\n")
}


// TODO: this is dogshit given by GPT, rewrite it. 
func PullFromNotion() {
    pageId, apiKey := GetTokens()
    // fmt.Println("📄 Page ID:", pageId)

    client := &http.Client{}
    endpoint := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", pageId)

    req, err := http.NewRequest("GET", endpoint, nil)
    checkErr(err)

    req.Header.Add("Notion-Version", "2022-06-28")
    req.Header.Add("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    checkErr(err)
    defer resp.Body.Close()

    bodyBytes, _ := io.ReadAll(resp.Body)
    if resp.StatusCode != http.StatusOK {
        fmt.Println("[NotionClient] Something's wrong with the dumb API client, please try again 10 mins later. If you wanna do it immediately, open notion :)")
        return
    }

	var notionList NotionListResponse
	err = json.Unmarshal(bodyBytes, &notionList)
	checkErr(err)

	for j := 0; j < len(notionList.Results); j++ {
		endpoint := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", notionList.Results[j].ID)

		req, err := http.NewRequest("GET", endpoint, nil)
		checkErr(err)

		req.Header.Add("Notion-Version", "2022-06-28")
		req.Header.Add("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		checkErr(err)
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			fmt.Println("[NotionClient] API error:", string(bodyBytes))
			return
		}

		var children NotionListResponse
		err = json.Unmarshal(bodyBytes, &children)
		checkErr(err)

		fmt.Println("Response for the page---> ")
		for _, block := range children.Results {
			for _, rt := range block.Paragraph.RichText {
				fmt.Println(rt.PlainText)
			}

			if block.ChildPage != nil {
				fmt.Println("Child page title:", block.ChildPage.Title)
			}
		}
	}


}