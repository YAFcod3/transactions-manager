package services

import (
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// AppServices centralizes all services in one struct
type AppServices struct {
	TransactionService     *ConversionService
	TransactionTypeService *TransactionTypeService
	// StatisticsService       *StatisticsService
	// SupportedCurrenciesService *SupportedCurrenciesService
}

// Initialize all services and return a centralized struct
func InitServices(db *mongo.Database, redisClient *redis.Client, codeGen *generate_transaction_code.CodeGenerator) *AppServices {
	return &AppServices{
		// TransactionService:      NewConversionService(db, redisClient, codeGen),
		TransactionTypeService: NewTransactionTypeService(db),
		// StatisticsService:       NewStatisticsService(db),
		// SupportedCurrenciesService: NewSupportedCurrenciesService(db),
	}
}
