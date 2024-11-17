package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupTransactionsHistoryRoutes(app *fiber.App, transactionsHistoryService *services.TransactionsHistoryService) {
	transactionsHistoryGroup := app.Group("/exchange/api/transactions")
	transactionsHistoryHandler := handlers.NewTransactionsHistoryHandler(transactionsHistoryService)
	transactionsHistoryGroup.Get("/", transactionsHistoryHandler.GetTransactionsHistory)

}
