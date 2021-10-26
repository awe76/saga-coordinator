package gateway

import (
	"log"

	"github.com/awe76/saga-coordinator/producer"
	"github.com/awe76/saga-coordinator/workflow"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	producer producer.Producer
}

type Handler interface {
	handleWorkflow() func(c *fiber.Ctx) error
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
