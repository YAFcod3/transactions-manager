package routes

import (
	"transactions-manager/app/services"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(apiGroup fiber.Router, codeGen *generate_transaction_code.CodeGenerator, services *services.AppServices) {
	SetupConversionRoutes(apiGroup, codeGen, services.TransactionTypeService)
	SetupSupportedCurrencyRoutes(apiGroup, services.SupportedCurrenciesService)
	SetupTransactionTypeRoutes(apiGroup, services.TransactionTypeService)
	SetupStatisticsRoutes(apiGroup, services.GetStatisticsService)
	SetupTransactionsHistoryRoutes(apiGroup, services.TransactionsHistoryService)

}
