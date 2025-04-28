package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/hasura/go-graphql-client"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/router"
)

func main() {
	// GraphQL Client
	graphQLClient := graphql.NewClient(fmt.Sprintf("https://%s/query", os.Getenv("API_URL")), nil)
	app := fiber.New(fiber.Config{
		AppName: "HYUabot-Kakao-Backend",
	})
	app.Use(logger.New())
	app.Use(func(ctx fiber.Ctx) error {
		ctx.Locals("graphQLClient", graphQLClient)
		return ctx.Next()
	})
	// Health Check
	app.Post("/healthcheck", router.GetHealthCheckMessage)
	// Listen on 3000
	err := app.Listen("0.0.0.0:3000", fiber.ListenConfig{
		EnablePrefork: true,
	})
	if err != nil {
		panic(err)
	}
}
