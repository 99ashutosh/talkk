package handlers

import "github.com/gofiber/fiber/v2"

type AppHandler struct{}

func NewAppHandler() *AppHandler {
	return &AppHandler{}
}

func (a *AppHandler) HandleGetIndex(ctx *fiber.Ctx) error {
	context := fiber.Map{
		"Title": "Talkk",
	}
	return ctx.Render("index", context)
}
