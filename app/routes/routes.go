package routes

import (
	"transactions-manager/app/services"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, codeGen *generate_transaction_code.CodeGenerator, services *services.AppServices) {
	SetupConversionRoutes(app, codeGen, services.SupportedCurrenciesService, services.TransactionTypeService)
	SetupSupportedCurrencyRoutes(app, services.SupportedCurrenciesService)
	SetupTransactionTypeRoutes(app, services.TransactionTypeService)
	SetupStatisticsRoutes(app, services.GetStatisticsService)
	SetupTransactionsHistoryRoutes(app, services.TransactionsHistoryService)

}
