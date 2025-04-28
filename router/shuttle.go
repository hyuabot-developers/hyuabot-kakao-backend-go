package router

import (
	"context"
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
	response := schema.SkillResponse{
		Version: "2.0",
		Template: schema.SkillTemplate{
			Outputs: []schema.Component{
				schema.SimpleText{Text: "API 서버 정상"},
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
