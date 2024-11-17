package handlers

import (
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
)

type SupportedCurrenciesHandler struct {
	Service *services.SupportedCurrenciesService
}

func NewSupportedCurrenciesHandler(service *services.SupportedCurrenciesService) *SupportedCurrenciesHandler {
	return &SupportedCurrenciesHandler{Service: service}
}

func (h *SupportedCurrenciesHandler) GetSupportedCurrencies(c *fiber.Ctx) error {
	currencies := h.Service.GetSupportedCurrencies()
	return c.JSON(fiber.Map{
		"supportedCurrencies": currencies,
	})
}
