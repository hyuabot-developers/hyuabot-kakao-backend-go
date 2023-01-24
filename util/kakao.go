package util

import "github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/model"

func SetResponse(template model.SkillTemplate) model.ServerResponse {
	return model.ServerResponse{Version: "2.0", Template: template}
}

func SetTemplate(outputs []model.Components, replies []model.QuickReply) model.SkillTemplate {
	return model.SkillTemplate{Outputs: outputs, QuickReplies: replies}
}

func setSimpleText(message string) model.SimpleTextResponse {
	return model.SimpleTextResponse{SimpleText: model.TextContent{message}}
}

func setBasicCard(title string, message string, buttons []model.CardButton) model.TextCard {
	return model.TextCard{Title: title, Description: message, Buttons: buttons}
}

func SetBasicCardCarousel(cards []model.TextCard) model.CarouselResponse {
	return model.CarouselResponse{Carousel: model.Carousel{Type: "basicCard", Items: cards}}
}
