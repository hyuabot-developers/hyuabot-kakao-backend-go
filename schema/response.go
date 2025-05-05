package schema

type (
	SkillResponse struct {
		Version  string        `json:"version"`
		Template SkillTemplate `json:"template"`
	}

	SkillTemplate struct {
		Outputs      []Component  `json:"outputs"`
		QuickReplies []QuickReply `json:"quickReplies"`
	}

	Component interface {
	}

	SimpleText struct {
		Component Component         `json:"component,omitempty"`
		Content   SimpleTextContent `json:"simpleText"`
	}

	SimpleTextContent struct {
		Text string `json:"text"`
	}

	SimpleImage struct {
		Component Component          `json:"component,omitempty"`
		Content   SimpleImageContent `json:"simpleImage"`
	}

	SimpleImageContent struct {
		ImageURL string `json:"imageUrl"`
		AltText  string `json:"altText"`
	}

	TextCard struct {
		Component Component       `json:"component,omitempty"`
		Content   TextCardContent `json:"textCard"`
	}

	TextCardContent struct {
		Title       string       `json:"title"`
		Description string       `json:"description"`
		Buttons     []CardButton `json:"buttons"`
	}

	BasicCard struct {
		Component Component        `json:"component,omitempty"`
		Content   BasicCardContent `json:"basicCard"`
	}

	BasicCardContent struct {
		Title       string             `json:"title"`
		Description string             `json:"description"`
		Thumbnail   BasicCardThumbnail `json:"thumbnail"`
		Buttons     []CardButton       `json:"buttons"`
	}

	BasicCardThumbnail struct {
		ImageURL string `json:"imageUrl"`
		AltText  string `json:"altText"`
	}

	CardButton struct {
		Label       string `json:"label"`
		Action      string `json:"action"`
		WebLinkURL  string `json:"webLinkUrl"`
		MessageText string `json:"messageText"`
		PhoneNumber string `json:"phoneNumber"`
		BlockID     string `json:"blockId"`
	}

	Carousel struct {
		Component Component       `json:"component,omitempty"`
		Content   CarouselContent `json:"carousel"`
	}

	CarouselContent struct {
		Type  string      `json:"type"`
		Items []Component `json:"items"`
	}

	QuickReply struct {
		Label       string            `json:"label"`
		Action      string            `json:"action"`
		MessageText string            `json:"messageText"`
		BlockID     string            `json:"blockId"`
		Extra       map[string]string `json:"extra"`
	}
)
