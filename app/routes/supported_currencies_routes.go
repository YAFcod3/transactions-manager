package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupSupportedCurrencyRoutes(app *fiber.App, service *services.SupportedCurrenciesService) {

	handler := handlers.NewSupportedCurrenciesHandler(service)

	group := app.Group("/exchange/api/currencies")
	group.Get("/", handler.GetSupportedCurrencies)
}
