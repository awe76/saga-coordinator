package gateway

import (
	"encoding/json"
	"log"

	"github.com/awe76/saga-coordinator/producer"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	producer producer.Producer
}

type Handler interface {
	createComment() func(c *fiber.Ctx) error
}

// createComment handler
func (h *handler) createComment() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Instantiate new Message struct
		cmt := new(Comment)

		//  Parse body into comment struct
		if err := c.BodyParser(cmt); err != nil {
			log.Println(err)
			c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": err,
			})
			return err
		}
		// convert body into bytes and send it to kafka
		cmtInBytes, err := json.Marshal(cmt)

		if err != nil {
			return err
		}
		h.producer.SendMessage("comments", cmtInBytes)

		// Return Comment in JSON format
		err = c.JSON(&fiber.Map{
			"success": true,
			"message": "Comment pushed successfully",
			"comment": cmt,
		})
		if err != nil {
			c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": "Error creating product",
			})
			return err
		}

		return err
	}
}
