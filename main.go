package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/backend"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/router"
)

const (
	backendRequestTimeout = 4 * time.Second
	defaultGraphQLURL     = "https://backend.hyuabot.app/graphql"
)

func main() {
	graphQLURL := os.Getenv("GRAPHQL_URL")
	if graphQLURL == "" {
		graphQLURL = defaultGraphQLURL
	}
	graphQLClient := backend.NewClient(graphQLURL, &http.Client{Timeout: backendRequestTimeout})
	app := fiber.New(fiber.Config{
		AppName: "HYUabot-Kakao-Backend",
	})
	app.Use(logger.New())
	app.Use(func(ctx fiber.Ctx) error {
		ctx.Locals(router.BackendClientLocal, graphQLClient)
		return ctx.Next()
	})
	// Routes
	app.Get("/healthcheck", router.GetHealthCheckMessage)
	app.Post("/healthcheck", router.GetHealthCheckMessage)
	app.Post("/shuttle", router.GetShuttleMessage)
	app.Post("/bus", router.GetBusMessage)
	app.Post("/cafeteria", router.GetCafeteriaMessage)
	app.Post("/subway", router.GetSubwayMessage)
	// Listen on the service port.
	err := app.Listen("0.0.0.0:38001", fiber.ListenConfig{
		EnablePrefork: false,
	})
	if err != nil {
		panic(err)
	}
}
