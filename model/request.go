package model

type UserMessage struct {
	Request UserRequest `json:"userRequest"`
}

type UserRequest struct {
	Message string `json:"utterance"`
}
