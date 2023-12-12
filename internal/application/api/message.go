package api

import (
	"devteambot/internal/domain/message"

	"github.com/gofiber/fiber/v2"
)

type MessageAPI interface {
	SendStandardMessage(*fiber.Ctx) error
}

type MessageHandler struct {
	Service message.Service `inject:"messageService"`
}

func (h *MessageHandler) SendStandardMessage(c *fiber.Ctx) error {
	err := h.Service.SendStandardMessage(c.Context(), c.Params("message"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to send standard message",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success send standard message",
	})
}
