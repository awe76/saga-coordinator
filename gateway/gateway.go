package gateway

import (
	"github.com/gofiber/fiber/v2"
)

// Comment struct
type Comment struct {
	Text string `form:"text" json:"text"`
}

type gateway struct {
	handler Handler
}

type Gateway interface {
	Run()
}

func (gw *gateway) Run() {
	app := fiber.New()
	api := app.Group("/api/v1") // /api

	createComment := gw.handler.createComment()
	api.Post("/comments", createComment)
	app.Listen(":3000")
}
