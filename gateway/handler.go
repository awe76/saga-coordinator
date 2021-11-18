package gateway

import (
	"log"

	"github.com/awe76/saga-coordinator/portfolio"
	"github.com/awe76/saga-coordinator/producer"
	"github.com/awe76/saga-coordinator/workflow"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	producer producer.Producer
}

type Handler interface {
	handleWorkflow() func(c *fiber.Ctx) error
	handleCreatePortfolio() func(c *fiber.Ctx) error
}

func (h *handler) handleWorkflow() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		w := new(workflow.Workflow)

		if err := c.BodyParser(w); err != nil {
			log.Println(err)
			c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": err,
			})
			return err
		}

		h.producer.SendMessage(workflow.WORKFLOW_START, w)

		err := c.JSON(&fiber.Map{
			"success":  true,
			"message":  "workflow pushed successfully",
			"workflow": w,
		})
		if err != nil {
			c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": "Error pushing workflow",
			})
			return err
		}

		return err
	}
}

func (h *handler) handleCreatePortfolio() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		payload := new(portfolio.CreateBuildingMapPortfolio)

		if err := c.BodyParser(payload); err != nil {
			log.Println(err)
			c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": err,
			})
			return err
		}

		w := &workflow.Workflow{
			Name:  "portfolio",
			Start: "s1",
			End:   "s3",
			Operations: []workflow.Operation{
				{
					From: "s1",
					To:   "s2",
					Name: "create-building-map-portfolio",
				},
				{
					From: "s2",
					To:   "s3",
					Name: "create-portfolio",
				},
			},
			Payload: payload,
		}

		h.producer.SendMessage(workflow.WORKFLOW_START, w)

		err := c.JSON(&fiber.Map{
			"success":  true,
			"message":  "create portfolio workflow pushed successfully",
			"workflow": w,
		})
		if err != nil {
			c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": "Error pushing create portfolio workflow",
			})
			return err
		}

		return err
	}
}
