package services

import (
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppServices struct {
	TransactionService         *ConversionService
	TransactionTypeService     *TransactionTypeService
	GetStatisticsService       *StatisticsService
	SupportedCurrenciesService *SupportedCurrenciesService
	TransactionsHistoryService *TransactionsHistoryService
}

func InitServices(db *mongo.Database, redisClient *redis.Client, codeGen *generate_transaction_code.CodeGenerator) *AppServices {
	return &AppServices{
		SupportedCurrenciesService: NewSupportedCurrenciesService(),
		TransactionTypeService:     NewTransactionTypeService(db),
		GetStatisticsService:       NewStatisticsService(db),
		TransactionsHistoryService: NewTransactionsHistoryService(db),
	}
}
