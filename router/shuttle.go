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

// Shuttle Stop ID.
const dormitoryStopID = "dormitory_o"
const shuttlecockOutStopID = "shuttlecock_o"
const stationStopID = "station"
const terminalStopID = "terminal"
const jungangStationStopID = "jungang_stn"
const shuttlecockInStopID = "shuttlecock_i"

// Shuttle Destination Group ID.
const stationDestination = "STATION"
const terminalDestination = "TERMINAL"
const jungangStationDestination = "JUNGANG"
const campusDestination = "CAMPUS"

// Header Text.
const headingStation = "한대앞\n"
const headingTerminal = "\n예술인\n"
const headingJungangStation = "\n중앙역\n"
const headingCampus = "캠퍼스\n"

// No Arrival Data.
const noArrivalText = "운행 없음\n"

type ShuttleTimetable struct {
	Tag         string
	Route       string
	Time        string
	Hour        int
	Minute      int
	Stop        string
	Destination string
}

func GenerateCardText(stopID string, resultMap map[string]map[string][]ShuttleTimetable) string {
	// Destination For Each Stop.
	destinationMap := map[string][]string{
		dormitoryStopID:      {stationDestination, terminalDestination, jungangStationDestination},
		shuttlecockOutStopID: {stationDestination, terminalDestination, jungangStationDestination},
		stationStopID:        {campusDestination, terminalDestination, jungangStationDestination},
		terminalStopID:       {campusDestination},
		jungangStationStopID: {campusDestination},
		shuttlecockInStopID:  {campusDestination},
	}
	headerMap := map[string]string{
		stationDestination:        headingStation,
		terminalDestination:       headingTerminal,
		jungangStationDestination: headingJungangStation,
		campusDestination:         headingCampus,
	}

	var cardText string
	for _, destination := range destinationMap[stopID] {
		cardText += headerMap[destination]
		result := resultMap[stopID][destination]
		if result == nil {
			cardText += noArrivalText
			continue
		}
		for _, timetable := range result {
			if timetable.Tag == "C" {
				cardText += fmt.Sprintf("순환 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			} else {
				cardText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			}
		}
	}
	return strings.Trim(cardText, "\n")
}

func QueryShuttleTimetable(ctx fiber.Ctx) []ShuttleTimetable {
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
		Shuttle struct {
			GroupedTimetable []ShuttleTimetable
		} `graphql:"shuttle(count: 2, startStr: $time, timestampStr: $timestamp, group: \"destination\")"`
	}
	variables := map[string]interface{}{
		"time":      currentTime.Format("15:04:03"),
		"timestamp": currentTime.Format("2006-01-02 15:04:03"),
	}
	queryError := client.Query(context.Background(), &query, variables)
	if queryError != nil {
		panic(queryError)
	}
	return query.Shuttle.GroupedTimetable
}

func GetShuttleMessage(ctx fiber.Ctx) error {
	body := new(schema.SkillPayload)
	if err := ctx.Bind().JSON(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Group shuttle timetable by stop and destination
	result := QueryShuttleTimetable(ctx)
	resultMap := make(map[string]map[string][]ShuttleTimetable)
	for _, timetable := range result {
		if resultMap[timetable.Stop] == nil {
			resultMap[timetable.Stop] = make(map[string][]ShuttleTimetable)
		}
		resultMap[timetable.Stop][timetable.Destination] = append(
			resultMap[timetable.Stop][timetable.Destination],
			timetable,
		)
	}
	// Create response text
	dormitoryText := GenerateCardText(dormitoryStopID, resultMap)
	shuttlecockOutText := GenerateCardText(shuttlecockOutStopID, resultMap)
	stationText := GenerateCardText(stationStopID, resultMap)
	terminalText := GenerateCardText(terminalStopID, resultMap)
	jungangStationText := GenerateCardText(jungangStationStopID, resultMap)
	shuttlecockInText := GenerateCardText(shuttlecockInStopID, resultMap)
	response := schema.SkillResponse{
		Version: "2.0",
		Template: schema.SkillTemplate{
			Outputs: []schema.Component{
				schema.Carousel{
					Content: schema.CarouselContent{
						Type: "textCard",
						Items: []schema.Content{
							schema.TextCardContent{
								Title:       "기숙사",
								Description: strings.Trim(dormitoryText, "\n"),
								Buttons:     []schema.CardButton{},
							},
							schema.TextCardContent{
								Title:       "셔틀콕",
								Description: strings.Trim(shuttlecockOutText, "\n"),
								Buttons:     []schema.CardButton{},
							},
							schema.TextCardContent{
								Title:       "한대앞",
								Description: strings.Trim(stationText, "\n"),
								Buttons:     []schema.CardButton{},
							},
							schema.TextCardContent{
								Title:       "예술인",
								Description: strings.Trim(terminalText, "\n"),
								Buttons:     []schema.CardButton{},
							},
							schema.TextCardContent{
								Title:       "중앙역",
								Description: strings.Trim(jungangStationText, "\n"),
								Buttons:     []schema.CardButton{},
							},
							schema.TextCardContent{
								Title:       "셔틀콕 건너편",
								Description: strings.Trim(shuttlecockInText, "\n"),
								Buttons:     []schema.CardButton{},
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
