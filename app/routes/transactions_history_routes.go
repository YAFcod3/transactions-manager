package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupTransactionsHistoryRoutes(apiGroup fiber.Router, transactionsHistoryService *services.TransactionsHistoryService) {
	transactionsHistoryHandler := handlers.NewTransactionsHistoryHandler(transactionsHistoryService)
	apiGroup.Get("/transactions", transactionsHistoryHandler.GetTransactionsHistory)

}
