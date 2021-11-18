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

	handleWorkflow := gw.handler.handleWorkflow()
	handleCreatePortfolio := gw.handler.handleCreatePortfolio()
	api.Post("/workflows", handleWorkflow)
	api.Post("/create-portfolio", handleCreatePortfolio)
	app.Listen(":3000")
}
