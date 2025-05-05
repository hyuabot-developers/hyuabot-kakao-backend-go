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

type SubwayStation struct {
	ID        string
	Realtime  SubwayRealtime
	Timetable SubwayTimetable
}

type SubwayRealtime struct {
	Up   []SubwayRealtimeItem
	Down []SubwayRealtimeItem
}

type SubwayRealtimeItem struct {
	Location string
	Time     float64
	Terminal SubwayTerminalStation
}

type SubwayTimetable struct {
	Up   []SubwayTimetableItem
	Down []SubwayTimetableItem
}

type SubwayTimetableItem struct {
	Time     string
	Terminal SubwayTerminalStation
}

type SubwayTerminalStation struct {
	Name string
}

func QuerySubwayDepartureData(ctx fiber.Ctx) []SubwayStation {
	// GraphQL Client and check API server status
	client, loaded := ctx.Locals("graphQLClient").(*graphql.Client)
	if !loaded {
		panic("GraphQL client not found")
	}
	// Get current datetime
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}
	currentTime := time.Now().In(location)
	// Query Shuttle Timetable
	var query struct {
		Subway []SubwayStation `graphql:"subway (id_: [\"K449\", \"K251\"], startStr: $start)"`
	}
	variables := map[string]interface{}{
		"start": currentTime.Format("15:04"),
	}
	queryError := client.Query(context.Background(), &query, variables)
	if queryError != nil {
		panic(queryError)
	}
	return query.Subway
}

func GenerateSubwayText(upHeaderText string, downHeaderText string, station SubwayStation) string {
	cardText := ""
	cardText += upHeaderText
	for index, realtime := range station.Realtime.Up {
		cardText += fmt.Sprintf("%s행 %d분 후 도착(%s)\n", realtime.Terminal.Name, int(realtime.Time), realtime.Location)
		if index == arrivalSectionLength-1 {
			break
		}
	}
	if len(station.Realtime.Up) < arrivalSectionLength {
		for index, timetable := range station.Timetable.Up {
			if index < arrivalSectionLength-len(station.Realtime.Up) {
				cardText += fmt.Sprintf("%s 출발\n", timetable.Time)
			}
		}
	}
	if len(station.Realtime.Up) == 0 && len(station.Timetable.Down) == 0 {
		cardText += "운행 정보 없음\n"
	}
	cardText += downHeaderText
	for index, realtime := range station.Realtime.Down {
		cardText += fmt.Sprintf("%s행 %d분 후 도착(%s)\n", realtime.Terminal.Name, int(realtime.Time), realtime.Location)
		if index == arrivalSectionLength-1 {
			break
		}
	}
	if len(station.Realtime.Down) < arrivalSectionLength {
		for index, timetable := range station.Timetable.Down {
			if index < arrivalSectionLength-len(station.Realtime.Down) {
				cardText += fmt.Sprintf("%s 출발\n", timetable.Time)
			}
		}
	}
	if len(station.Realtime.Down) == 0 && len(station.Timetable.Down) == 0 {
		cardText += "운행 정보 없음\n"
	}
	return cardText
}

func GetSubwayMessage(ctx fiber.Ctx) error {
	body := new(schema.SkillPayload)
	if err := ctx.Bind().JSON(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	result := QuerySubwayDepartureData(ctx)
	resultMap := make(map[string]SubwayStation)
	for _, station := range result {
		resultMap[station.ID] = station
	}
	// Create response text
	line4Text := GenerateSubwayText("당고개 방면\n", "\n오이도 방면\n", resultMap["K449"])
	lineSuinText := GenerateSubwayText("청량리 방면\n", "\n인천 방면\n", resultMap["K251"])
	response := schema.SkillResponse{
		Version: "2.0",
		Template: schema.SkillTemplate{
			Outputs: []schema.Component{
				schema.Carousel{
					Type: "textCard",
					Items: []schema.Component{
						schema.TextCard{
							Title:       "4호선",
							Description: strings.Trim(line4Text, "\n"),
							Buttons:     []schema.CardButton{},
						},
						schema.TextCard{
							Title:       "수인분당선",
							Description: strings.Trim(lineSuinText, "\n"),
							Buttons:     []schema.CardButton{},
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
