package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupTransactionTypeRoutes(apiGroup fiber.Router, transactionTypeService *services.TransactionTypeService) {
	transactionTypeGroup := apiGroup.Group("/settings/transactions-types")
	transactionTypeHandler := handlers.NewTransactionTypeHandler(transactionTypeService)

	transactionTypeGroup.Get("/", transactionTypeHandler.GetTransactionTypes)
	transactionTypeGroup.Post("/", transactionTypeHandler.CreateTransactionType)
	transactionTypeGroup.Get("/:id", transactionTypeHandler.GetTransactionTypeByID)
	transactionTypeGroup.Patch("/:id", transactionTypeHandler.UpdateTransactionType)
	transactionTypeGroup.Delete("/:id", transactionTypeHandler.DeleteTransactionType)
}
