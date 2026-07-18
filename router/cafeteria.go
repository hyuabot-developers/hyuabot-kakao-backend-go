package router

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/schema"
)

const (
	lunchStartHour      = 9
	lunchEndHour        = 17
	ericaCampusID       = 2
	maxCafeteriaCards   = 5
	maxMenusPerListCard = 4
	quickReplyCapacity  = 7
)

const cafeteriaQuery = `
query KakaoCafeteria($date: Date!, $campus: Int!) {
  cafeteria(input: {date: $date, campus: $campus}) {
    seq
    name
    menus { type food price }
  }
}`

type cafeteriaResult struct {
	Cafeteria []cafeteria `json:"cafeteria"`
}

type cafeteria struct {
	Seq   int    `json:"seq"`
	Name  string `json:"name"`
	Menus []menu `json:"menus"`
}

type menu struct {
	Type  string `json:"type"`
	Food  string `json:"food"`
	Price string `json:"price"`
}

func GetCafeteriaMessage(ctx fiber.Ctx) error {
	payload, err := bindSkillPayload(ctx)
	if err != nil {
		return badRequest(ctx, err)
	}
	client, err := backendClient(ctx)
	if err != nil {
		return backendFailure(ctx, "cafeteria", "학식", err)
	}

	currentTime := currentServiceTime()
	mealType := requestedMealType(payload, currentTime.Hour())
	queryCtx, cancel := queryContext()
	defer cancel()
	var result cafeteriaResult
	if queryErr := client.Query(queryCtx, cafeteriaQuery, map[string]any{
		"date":   currentTime.Format("2006-01-02"),
		"campus": ericaCampusID,
	}, &result); queryErr != nil {
		return backendFailure(ctx, "cafeteria", "학식", queryErr)
	}

	items := make([]any, 0, min(len(result.Cafeteria), maxCafeteriaCards))
	for _, cafeteria := range result.Cafeteria {
		if len(items) == maxCafeteriaCards {
			break
		}
		items = append(items, cafeteriaCard(cafeteria, mealType))
	}
	if len(items) == 0 {
		return ctx.JSON(schema.TextResponse(
			"오늘 등록된 학식 정보가 없어요.",
			cafeteriaQuickReplies(mealType),
		))
	}

	return ctx.JSON(schema.CarouselResponse("listCard", items, cafeteriaQuickReplies(mealType)))
}

func requestedMealType(payload *schema.SkillPayload, hour int) string {
	for _, candidate := range []string{payload.Action.Params["meal"], payload.Action.Params["type"]} {
		switch candidate {
		case "조식", "중식", "석식":
			return candidate
		}
	}
	if strings.Contains(payload.UserRequest.Utterance, "조식") {
		return "조식"
	}
	if strings.Contains(payload.UserRequest.Utterance, "석식") {
		return "석식"
	}
	if hour < lunchStartHour {
		return "조식"
	}
	if hour >= lunchEndHour {
		return "석식"
	}
	return "중식"
}

func cafeteriaCard(cafeteria cafeteria, mealType string) schema.ListCard {
	items := make([]schema.ListItem, 0, maxMenusPerListCard)
	for _, menu := range cafeteria.Menus {
		if menu.Type != mealType || len(items) == maxMenusPerListCard {
			continue
		}
		items = append(items, schema.ListItem{
			Title:       strings.TrimSpace(menu.Food),
			Description: strings.TrimSpace(menu.Price),
		})
	}
	if len(items) == 0 {
		items = append(items, schema.ListItem{Title: "등록된 메뉴가 없어요"})
	}
	return schema.ListCard{
		Header:  schema.ListItem{Title: cafeteria.Name + " · " + mealType},
		Items:   items,
		Buttons: []schema.Button{appButton("/cafeteria", "전체 메뉴")},
	}
}

func cafeteriaQuickReplies(currentMeal string) []schema.QuickReply {
	replies := make([]schema.QuickReply, 0, quickReplyCapacity)
	for _, mealType := range []string{"조식", "중식", "석식"} {
		if mealType == currentMeal {
			continue
		}
		replies = append(replies, schema.QuickReply{
			Label:       mealType,
			Action:      "message",
			MessageText: mealType + " 메뉴",
		})
	}
	replies = append(replies, navigationQuickReplies("학식")...)
	return replies
}
