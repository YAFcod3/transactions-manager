package handlers

import (
	"time"
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
)

type TransactionsHistoryHandler struct {
	Service *services.TransactionsHistoryService
}

func NewTransactionsHistoryHandler(service *services.TransactionsHistoryService) *TransactionsHistoryHandler {
	return &TransactionsHistoryHandler{Service: service}
}

func (h *TransactionsHistoryHandler) GetTransactionsHistory(c *fiber.Ctx) error {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	transactionType := c.Query("transactionType")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_DATE_FORMAT",
				"message": "Invalid startDate format. Use RFC3339.",
			})
		}
	}
	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_DATE_FORMAT",
				"message": "Invalid endDate format. Use RFC3339.",
			})
		}
	}

	if startDate.IsZero() {
		startDate = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	history, err := h.Service.GetTransactionsHistory(c.Context(), startDate, endDate, transactionType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "TRANSACTION_HISTORY_ERROR",
			"message": "Failed to fetch transactions history.",
		})
	}

	return c.JSON(fiber.Map{
		"data": history,
	})
}
