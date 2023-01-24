package model

type ServerResponse struct {
	Version  string        `json:"version"`
	Template SkillTemplate `json:"template"`
}

type SkillTemplate struct {
	Outputs      []Components `json:"outputs"`
	QuickReplies []QuickReply `json:"quickReplies"`
}

type Components interface{}

type SimpleTextResponse struct {
	SimpleText TextContent `json:"simpleText"`
}

type TextContent struct {
	Text string `json:"text"`
}

type BasicCardResponse struct {
	Card TextCard `json:"basicCard"`
}

type TextCard struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Buttons     []CardButton `json:"buttons"`
}

type CarouselResponse struct {
	Carousel Carousel `json:"carousel"`
}

type Carousel struct {
	Type  string     `json:"type"`
	Items []TextCard `json:"items"`
}

type QuickReply struct {
	Action      string `json:"action"`
	Label       string `json:"label"`
	MessageText string `json:"messageText"`
	BlockID     string `json:"blockId"`
}

type CardButton struct {
	Action string `json:"action"`
	Label  string `json:"label"`
	Link   string `json:"webLinkUrl"`
}
