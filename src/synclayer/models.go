package synclayer

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