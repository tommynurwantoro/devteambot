package api

import (
	"devteambot/internal/domain/sholat"

	"github.com/gofiber/fiber/v2"
)

type SholatAPI interface {
	GetSholatSchedule(*fiber.Ctx) error
	SendReminderSholat(*fiber.Ctx) error
}

type SholatHandler struct {
	Service sholat.Service `inject:"sholatService"`
}

func (h *SholatHandler) GetSholatSchedule(c *fiber.Ctx) error {
	err := h.Service.GetTodaySchedule(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get sholat schedule",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success get sholat schedule",
	})
}

func (h *SholatHandler) SendReminderSholat(c *fiber.Ctx) error {
	err := h.Service.SendReminder(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to send reminder sholat",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success send reminder sholat",
	})
}
