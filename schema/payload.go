package schema

type (
	SkillPayload struct {
		Bot    Bot    `json:"bot"`
		Intent Intent `json:"intent"`
	}

	Bot struct {
		BotID string `json:"id"`
		Name  string `json:"name"`
	}

	Intent struct {
		BlockID string      `json:"id"`
		Name    string      `json:"name"`
		Extra   IntentExtra `json:"extra"`
	}

	IntentExtra struct {
		Knowledge Knowledge `json:"knowledge"`
	}

	Knowledge struct {
		ResponseType     string             `json:"responseType"`
		MatchedKnowledge []MatchedKnowledge `json:"matchedKnowledges"`
	}

	MatchedKnowledge struct {
		Categories []string `json:"categories"`
		Question   string   `json:"question"`
		Answer     string   `json:"answer"`
		LandingURL string   `json:"landingUrl"`
		ImageURL   string   `json:"imageUrl"`
	}

	Action struct {
		SkillID string         `json:"id"`
		Name    string         `json:"name"`
		Extra   map[string]any `json:"clientExtra"`
	}

	UserRequest struct {
		Timezone  string `json:"timezone"`
		Block     Block  `json:"block"`
		Utterance string `json:"utterance"`
		Language  string `json:"lang"`
		User      User   `json:"user"`
	}

	Block struct {
		BlockID string `json:"id"`
		Name    string `json:"name"`
	}

	User struct {
		UserID     string         `json:"id"`
		Type       string         `json:"type"`
		Properties UserProperties `json:"properties"`
	}

	UserProperties struct {
		PlusFriendUserKey string `json:"plusFriendUserKey"`
		AppUserID         string `json:"appUserId"`
		IsFriend          bool   `json:"isFriend"`
	}
)
