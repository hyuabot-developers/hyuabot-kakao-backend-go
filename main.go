package main

import "github.com/gofiber/fiber/v3"

func main() {
	app := fiber.New(fiber.Config{
		AppName: "HYUabot-Kakao-Backend",
	})
	// Health Check
	app.Get("/health", func(ctx fiber.Ctx) error {
		return ctx.SendString("OK")
	})
	// Listen on 3000
	err := app.Listen(":3000", fiber.ListenConfig{
		EnablePrefork: true,
	})
	if err != nil {
		panic(err)
	}
}
