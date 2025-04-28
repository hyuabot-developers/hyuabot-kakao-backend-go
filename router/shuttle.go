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

func GetShuttleMessage(ctx fiber.Ctx) error {
	body := new(schema.SkillPayload)
	if err := ctx.Bind().JSON(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// GraphQL Client and check API server status
	client, loaded := ctx.Locals("graphQLClient").(*graphql.Client)
	if !loaded {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "GraphQL client not found",
		})
	}
	// Get current datetime
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	time.Local = location
	currentTime := time.Now()
	// Query Shuttle Timetable
	type DateTime time.Time
	type Time time.Time
	var query struct {
		Shuttle struct {
			GroupedTimetable []struct {
				Tag         string
				Route       string
				Time        string
				Hour        int
				Minute      int
				Stop        string
				Destination string
			}
		} `graphql:"shuttle(count: 2, startStr: $time, timestampStr: $timestamp, group: \"destination\")"`
	}
	variables := map[string]interface{}{
		"time":      currentTime.Format("15:04:03"),
		"timestamp": currentTime.Format("2006-01-02 15:04:03"),
	}
	queryError := client.Query(context.Background(), &query, variables)
	if queryError != nil {
		response := schema.SkillResponse{
			Version: "2.0",
			Template: schema.SkillTemplate{
				Outputs: []schema.Component{
					schema.SimpleText{Text: queryError.Error()},
				},
				QuickReplies: []schema.QuickReply{},
			},
		}
		return ctx.JSON(response)
	}
	// Create dormitory response text
	var dormitoryText string
	var count = 0
	dormitoryText += "한대앞\n"
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "dormitory_o" && timetable.Destination == "STATION" {
			count += 1
			if timetable.Tag == "C" {
				dormitoryText += fmt.Sprintf("순환 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			} else {
				dormitoryText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			}
		}
	}
	if count == 0 {
		dormitoryText += "운행 없음\n"
	}
	dormitoryText += "\n예술인\n"
	count = 0
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "dormitory_o" && timetable.Destination == "TERMINAL" {
			count += 1
			if timetable.Tag == "C" {
				dormitoryText += fmt.Sprintf("순환 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			} else {
				dormitoryText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			}
		}
	}
	if count == 0 {
		dormitoryText += "운행 없음\n"
	}
	dormitoryText += "\n중앙역\n"
	count = 0
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "dormitory_o" && timetable.Destination == "JUNGANG" {
			count += 1
			dormitoryText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
		}
	}
	if count == 0 {
		dormitoryText += "운행 없음\n"
	}
	// Create shuttlecock out response text
	var shuttlecockOutText string
	count = 0
	shuttlecockOutText += "한대앞\n"
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "shuttlecock_o" && timetable.Destination == "STATION" {
			count += 1
			if timetable.Tag == "C" {
				shuttlecockOutText += fmt.Sprintf("순환 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			} else {
				shuttlecockOutText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			}
		}
	}
	if count == 0 {
		shuttlecockOutText += "운행 없음\n"
	}
	shuttlecockOutText += "\n예술인\n"
	count = 0
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "shuttlecock_o" && timetable.Destination == "TERMINAL" {
			count += 1
			if timetable.Tag == "C" {
				shuttlecockOutText += fmt.Sprintf("순환 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			} else {
				shuttlecockOutText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			}
		}
	}
	if count == 0 {
		shuttlecockOutText += "운행 없음\n"
	}
	shuttlecockOutText += "\n중앙역\n"
	count = 0
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "shuttlecock_o" && timetable.Destination == "JUNGANG" {
			count += 1
			shuttlecockOutText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
		}
	}
	if count == 0 {
		shuttlecockOutText += "운행 없음\n"
	}
	// Create station response text
	var stationText string
	count = 0
	stationText += "캠퍼스\n"
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "station" && timetable.Destination == "CAMPUS" {
			count += 1
			if timetable.Tag == "C" {
				stationText += fmt.Sprintf("순환 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			} else {
				stationText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
			}
		}
	}
	if count == 0 {
		stationText += "운행 없음\n"
	}
	stationText += "\n예술인\n"
	count = 0
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "station" && timetable.Destination == "TERMINAL" {
			count += 1
			stationText += fmt.Sprintf("순환 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
		}
	}
	if count == 0 {
		stationText += "운행 없음\n"
	}
	stationText += "\n중앙역\n"
	count = 0
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "station" && timetable.Destination == "JUNGANG" {
			count += 1
			stationText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
		}
	}
	if count == 0 {
		stationText += "운행 없음\n"
	}
	// Create terminal response text
	var terminalText string
	count = 0
	terminalText += "캠퍼스\n"
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "terminal" {
			count += 1
			terminalText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
		}
	}
	if count == 0 {
		terminalText += "운행 없음\n"
	}
	// Create jungang station response text
	var jungangStationText string
	count = 0
	jungangStationText += "캠퍼스\n"
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "jungang_stn" {
			count += 1
			jungangStationText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
		}
	}
	if count == 0 {
		jungangStationText += "운행 없음\n"
	}
	// Create shuttlecock in response text
	var shuttlecockInText string
	count = 0
	shuttlecockInText += "기숙사\n"
	for _, timetable := range query.Shuttle.GroupedTimetable {
		if timetable.Stop == "shuttlecock_i" {
			count += 1
			shuttlecockInText += fmt.Sprintf("직행 %02d:%02d 출발\n", timetable.Hour, timetable.Minute)
		}
	}
	if count == 0 {
		shuttlecockInText += "운행 없음\n"
	}
	response := schema.SkillResponse{
		Version: "2.0",
		Template: schema.SkillTemplate{
			Outputs: []schema.Component{
				schema.Carousel{
					Type: "textCard",
					Items: []schema.Component{
						schema.TextCard{
							Title:       "기숙사",
							Description: strings.Trim(dormitoryText, "\n"),
							Buttons:     []schema.CardButton{},
						},
						schema.TextCard{
							Title:       "셔틀콕",
							Description: strings.Trim(shuttlecockOutText, "\n"),
							Buttons:     []schema.CardButton{},
						},
						schema.TextCard{
							Title:       "한대앞",
							Description: strings.Trim(stationText, "\n"),
							Buttons:     []schema.CardButton{},
						},
						schema.TextCard{
							Title:       "예술인",
							Description: strings.Trim(terminalText, "\n"),
							Buttons:     []schema.CardButton{},
						},
						schema.TextCard{
							Title:       "중앙역",
							Description: strings.Trim(jungangStationText, "\n"),
							Buttons:     []schema.CardButton{},
						},
						schema.TextCard{
							Title:       "셔틀콕 건너편",
							Description: strings.Trim(shuttlecockInText, "\n"),
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
