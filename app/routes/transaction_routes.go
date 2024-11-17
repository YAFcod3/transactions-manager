package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

func SetupConversionRoutes(app *fiber.App, codeGen *generate_transaction_code.CodeGenerator) {
	transactionHandler := handlers.NewTransactionHandler(codeGen)
	app.Post("/api/conversion", transactionHandler.HandleTransaction)
}
