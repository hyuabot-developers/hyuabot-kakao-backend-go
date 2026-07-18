package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/schema"
)

const healthcheckQuery = `query Healthcheck { __typename }`

func GetHealthCheckMessage(ctx fiber.Ctx) error {
	client, err := backendClient(ctx)
	if err != nil {
		return healthcheckFailure(ctx, err)
	}

	queryCtx, cancel := queryContext()
	defer cancel()
	var result struct {
		TypeName string `json:"__typename"`
	}
	if queryErr := client.Query(queryCtx, healthcheckQuery, nil, &result); queryErr != nil {
		return healthcheckFailure(ctx, queryErr)
	}

	if ctx.Method() == fiber.MethodGet {
		return ctx.JSON(fiber.Map{"status": "ok"})
	}
	return ctx.JSON(schema.TextResponse("API 서버가 정상적으로 동작하고 있어요.", nil))
}

func healthcheckFailure(ctx fiber.Ctx, err error) error {
	if ctx.Method() == fiber.MethodGet {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "unavailable"})
	}
	return backendFailure(ctx, "healthcheck", "상태 확인", err)
}
