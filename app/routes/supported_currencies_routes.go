package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupSupportedCurrencyRoutes(apiGroup fiber.Router, service *services.SupportedCurrenciesService) {

	handler := handlers.NewSupportedCurrenciesHandler(service)

	apiGroup.Get("/currencies", handler.GetSupportedCurrencies)
}
