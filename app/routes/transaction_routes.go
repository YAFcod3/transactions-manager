package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/middleware"
	"transactions-manager/app/services"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

func SetupConversionRoutes(app *fiber.App, codeGen *generate_transaction_code.CodeGenerator, transactionTypeService *services.TransactionTypeService) {
	transactionHandler := handlers.NewTransactionHandler(codeGen, transactionTypeService)
	duplicateMiddleware := middleware.VerifyTransactionDuplicated()

	app.Post("/exchange/api/conversion", duplicateMiddleware, transactionHandler.HandleTransaction)
}
