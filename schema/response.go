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

	QuickReply struct {
		Label       string `json:"label"`
		Action      string `json:"action"`
		MessageText string `json:"messageText"`
		BlockID     string `json:"blockId"`
		Extra       any    `json:"extra"`
	}
)
