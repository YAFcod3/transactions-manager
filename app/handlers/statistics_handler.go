package handlers

import (
	"net/http"
	"time"
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
)

type StatisticsHandler struct {
	Service *services.StatisticsService
}

func NewStatisticsHandler(service *services.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{Service: service}
}

func (h *StatisticsHandler) GetStatisticsHandler(c *fiber.Ctx) error {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	var start, end time.Time
	var err error

	if startDate != "" {
		start, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_DATE_FORMAT",
				"message": "Start date is invalid. Use RFC3339 format.",
			})
		}
	} else {
		start = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	if endDate != "" {
		end, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_DATE_FORMAT",
				"message": "End date is invalid. Use RFC3339 format.",
			})
		}
	} else {
		end = time.Now()
	}
	stats, err := h.Service.GetStatisticsService(start, end)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve statistics",
		})
	}

	return c.JSON(stats)
}
