package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hyuabot-developers/hyuabot-kakao-i-server-golang/controller"
)

func main() {
	app := fiber.New()
	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://*.hyuabot.app, http://localhost:8100, http://192.168.*.*:8100",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Post("/bus/arrival", controller.BusArrival)
	app.Post("/shuttle/arrival", controller.ShuttleArrival)
	app.Post("/subway/arrival", controller.Subway)
	app.Post("/cafeteria/menu", controller.CafeteriaMenu)

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
