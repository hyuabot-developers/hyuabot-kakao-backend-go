package schema

const Version = "2.0"

type SkillResponse struct {
	Version  string        `json:"version"`
	Template SkillTemplate `json:"template"`
}

type SkillTemplate struct {
	Outputs      []any        `json:"outputs"`
	QuickReplies []QuickReply `json:"quickReplies,omitempty"`
}

type SimpleTextOutput struct {
	SimpleText SimpleText `json:"simpleText"`
}

type SimpleText struct {
	Text string `json:"text"`
}

type CarouselOutput struct {
	Carousel Carousel `json:"carousel"`
}

type Carousel struct {
	Type  string `json:"type"`
	Items []any  `json:"items"`
}

type ListCard struct {
	Header  ListItem   `json:"header"`
	Items   []ListItem `json:"items"`
	Buttons []Button   `json:"buttons,omitempty"`
}

type ListItem struct {
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	ImageURL    string         `json:"imageUrl,omitempty"`
	Action      string         `json:"action,omitempty"`
	MessageText string         `json:"messageText,omitempty"`
	BlockID     string         `json:"blockId,omitempty"`
	Extra       map[string]any `json:"extra,omitempty"`
}

type ItemCard struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	ItemList    []ItemList `json:"itemList"`
	Buttons     []Button   `json:"buttons,omitempty"`
}

type ItemList struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Button struct {
	Label       string         `json:"label"`
	Action      string         `json:"action"`
	WebLinkURL  string         `json:"webLinkUrl,omitempty"`
	MessageText string         `json:"messageText,omitempty"`
	PhoneNumber string         `json:"phoneNumber,omitempty"`
	BlockID     string         `json:"blockId,omitempty"`
	Extra       map[string]any `json:"extra,omitempty"`
}

type QuickReply struct {
	Label       string         `json:"label"`
	Action      string         `json:"action"`
	MessageText string         `json:"messageText,omitempty"`
	BlockID     string         `json:"blockId,omitempty"`
	Extra       map[string]any `json:"extra,omitempty"`
}

func TextResponse(text string, quickReplies []QuickReply) SkillResponse {
	return SkillResponse{
		Version: Version,
		Template: SkillTemplate{
			Outputs:      []any{SimpleTextOutput{SimpleText: SimpleText{Text: text}}},
			QuickReplies: quickReplies,
		},
	}
}

func CarouselResponse(carouselType string, items []any, quickReplies []QuickReply) SkillResponse {
	return SkillResponse{
		Version: Version,
		Template: SkillTemplate{
			Outputs: []any{
				CarouselOutput{Carousel: Carousel{Type: carouselType, Items: items}},
			},
			QuickReplies: quickReplies,
		},
	}
}
