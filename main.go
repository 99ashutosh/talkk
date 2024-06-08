package main

import (
	"github.com/99ashutosh/talkk/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	engine := html.New("./views", ".html")
	appHandler := handlers.NewAppHandler()

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static/", "./static")
	server := NewWebSocket()

	app.Get("/", appHandler.HandleGetIndex)
	app.Get("/ws", websocket.New((func(ctx *websocket.Conn) {
		server.HandleWebSocket(ctx)
	})))

	go server.HandleMessages()

	app.Get("/metrics", monitor.New(monitor.Config{Title: "Talkk monitor"}))

	app.Listen(":3000")
}
