package router

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/hasura/go-graphql-client"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/schema"
)

const lunchStartHour = 9
const lunchEndHour = 17

type Cafeteria struct {
	ID   int
	Menu []Menu
}

type Menu struct {
	Menu  string
	Price string
}

func QueryCafeteriaDepartureData(ctx fiber.Ctx, date string, feedType string) []Cafeteria {
	// GraphQL Client and check API server status
	client, loaded := ctx.Locals("graphQLClient").(*graphql.Client)
	if !loaded {
		panic("GraphQL client not found")
	}
	// Query cafeteria menu
	var query struct {
		Menu []Cafeteria `graphql:"menu (dateStr: $dateStr, campusId: 2, type_: [$type_])"`
	}
	variables := map[string]interface{}{
		"dateStr": date,
		"type_":   feedType,
	}
	queryError := client.Query(context.Background(), &query, variables)
	if queryError != nil {
		panic(queryError)
	}
	return query.Menu
}

func GenerateCafeteriaText(cafeteria Cafeteria) string {
	cardText := ""
	for _, menu := range cafeteria.Menu {
		cardText += fmt.Sprintf("%s\n%s\n", strings.TrimSpace(menu.Menu), menu.Price)
	}
	return cardText
}

func GetCafeteriaMessage(ctx fiber.Ctx) error {
	body := new(schema.SkillPayload)
	if err := ctx.Bind().JSON(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Get current datetime
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}
	currentTime := time.Now().In(location)
	// Set food type
	feedType := "중식"
	if currentTime.Hour() >= lunchEndHour {
		feedType = "석식"
	} else if currentTime.Hour() < lunchStartHour {
		feedType = "조식"
	}
	// Group shuttle timetable by stop and destination
	result := QueryCafeteriaDepartureData(ctx, time.Now().Format("2006-01-02"), feedType)
	resultMap := make(map[int]Cafeteria)
	for _, cafeteria := range result {
		resultMap[cafeteria.ID] = cafeteria
	}
	staffCafeteriaText := GenerateCafeteriaText(resultMap[11])
	studentCafeteriaText := GenerateCafeteriaText(resultMap[12])
	dormitoryCafeteriaText := GenerateCafeteriaText(resultMap[13])
	foodCourtCafeteriaText := GenerateCafeteriaText(resultMap[14])
	businessCafeteriaText := GenerateCafeteriaText(resultMap[15])
	response := schema.SkillResponse{
		Version: "2.0",
		Template: schema.SkillTemplate{
			Outputs: []schema.Component{
				schema.Carousel{
					Content: schema.CarouselContent{
						Type: "textCard",
						Items: []schema.Component{
							schema.TextCard{
								Content: schema.TextCardContent{
									Title:       fmt.Sprintf("%s(%s)", "교직원식당", feedType),
									Description: strings.Trim(staffCafeteriaText, "\n"),
									Buttons:     []schema.CardButton{},
								},
							},
							schema.TextCard{
								Content: schema.TextCardContent{
									Title:       fmt.Sprintf("%s(%s)", "학생식당", feedType),
									Description: strings.Trim(studentCafeteriaText, "\n"),
									Buttons:     []schema.CardButton{},
								},
							},
							schema.TextCard{
								Content: schema.TextCardContent{
									Title:       fmt.Sprintf("%s(%s)", "창의인재원식당", feedType),
									Description: strings.Trim(dormitoryCafeteriaText, "\n"),
									Buttons:     []schema.CardButton{},
								},
							},
							schema.TextCard{
								Content: schema.TextCardContent{
									Title:       fmt.Sprintf("%s(%s)", "푸드코트", feedType),
									Description: strings.Trim(foodCourtCafeteriaText, "\n"),
									Buttons:     []schema.CardButton{},
								},
							},
							schema.TextCard{
								Content: schema.TextCardContent{
									Title:       fmt.Sprintf("%s(%s)", "창업보육센터", feedType),
									Description: strings.Trim(businessCafeteriaText, "\n"),
									Buttons:     []schema.CardButton{},
								},
							},
						},
					},
				},
			},
			QuickReplies: []schema.QuickReply{
				{
					Label:       "휴아봇 앱 설치",
					Action:      "block",
					MessageText: "휴아봇 앱 설치",
					BlockID:     "6077ca2de2039a2ba38c755f",
					Extra:       map[string]string{},
				},
			},
		},
	}
	return ctx.JSON(response)
}
