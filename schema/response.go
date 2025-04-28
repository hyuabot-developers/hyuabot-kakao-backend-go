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
		Component Component `json:"component,omitempty"`
		Text      string    `json:"text"`
	}

	SimpleImage struct {
		Component Component `json:"component,omitempty"`
		ImageURL  string    `json:"imageUrl"`
		AltText   string    `json:"altText"`
	}

	TextCard struct {
		Component   Component    `json:"component,omitempty"`
		Title       string       `json:"title"`
		Description string       `json:"description"`
		Buttons     []CardButton `json:"buttons"`
	}

	BasicCard struct {
		Component   Component          `json:"component,omitempty"`
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
		Component Component   `json:"component,omitempty"`
		Type      string      `json:"type"`
		Items     []Component `json:"items"`
	}

	QuickReply struct {
		Label       string `json:"label"`
		Action      string `json:"action"`
		MessageText string `json:"messageText"`
		BlockID     string `json:"blockId"`
		Extra       any    `json:"extra"`
	}
)
