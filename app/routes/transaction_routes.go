package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/services"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

func SetupConversionRoutes(app *fiber.App, codeGen *generate_transaction_code.CodeGenerator, supportedCurrenciesService *services.SupportedCurrenciesService) {
	transactionHandler := handlers.NewTransactionHandler(codeGen, supportedCurrenciesService)
	app.Post("/api/conversion", transactionHandler.HandleTransaction)
}
