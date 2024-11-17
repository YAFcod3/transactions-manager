package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/services"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

func SetupConversionRoutes(app *fiber.App, codeGen *generate_transaction_code.CodeGenerator, supportedCurrenciesService *services.SupportedCurrenciesService, transactionTypeService *services.TransactionTypeService) {
	transactionHandler := handlers.NewTransactionHandler(codeGen, supportedCurrenciesService, transactionTypeService)
	app.Post("/exchange/api/conversion", transactionHandler.HandleTransaction)
}
