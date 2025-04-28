package router

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/hasura/go-graphql-client"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/schema"
)

func GetHealthCheckMessage(ctx fiber.Ctx) error {
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
	var query struct {
		Health bool
	}
	queryError := client.Query(context.Background(), &query, nil)
	if queryError != nil || !query.Health {
		response := schema.SkillResponse{
			Version: "2.0",
			Template: schema.SkillTemplate{
				Outputs: []schema.Component{
					schema.SimpleText{Text: "API 서버 비정상"},
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
			QuickReplies: []schema.QuickReply{},
		},
	}
	return ctx.JSON(response)
}
