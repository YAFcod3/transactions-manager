package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"transactions-manager/app/database"
	"transactions-manager/app/models"
	"transactions-manager/app/utils"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversionService struct {
	CodeGen                    *generate_transaction_code.CodeGenerator
	SupportedCurrenciesService *SupportedCurrenciesService
	TransactionTypeService     *TransactionTypeService
	SupportedCurrencies        []string
}

func NewConversionService(codeGen *generate_transaction_code.CodeGenerator, transactionTypeService *TransactionTypeService) *ConversionService {
	currencies := os.Getenv("SUPPORTED_CURRENCIES")
	if currencies == "" {
		currencies = "USD,EUR,GBP,JPY,CAD,AUD"
	}

	supportedCurrencies := []string{}
	for _, currency := range strings.Split(currencies, ",") {
		supportedCurrencies = append(supportedCurrencies, strings.TrimSpace(currency))
	}

	return &ConversionService{
		CodeGen:                codeGen,
		TransactionTypeService: transactionTypeService,
		SupportedCurrencies:    supportedCurrencies,
	}
}

func (s *ConversionService) ProcessTransaction(req models.TransactionRequest, userID string) (map[string]interface{}, error) {
	if !s.isCurrencySupported(req.FromCurrency) {
		supportedCurrencies := strings.Join(s.SupportedCurrencies, ", ")
		return nil, fmt.Errorf("INVALID_CURRENCY: The currency '%s' is not supported. Please use one of the supported currencies: %v", req.FromCurrency, supportedCurrencies)
	}
	if !s.isCurrencySupported(req.ToCurrency) {
		supportedCurrencies := strings.Join(s.SupportedCurrencies, ", ")
		return nil, fmt.Errorf("INVALID_CURRENCY: The currency '%s' is not supported. Please use one of the supported currencies: %v", req.ToCurrency, supportedCurrencies)
	}

	transactionTypeID, err := primitive.ObjectIDFromHex(req.TransactionType)
	if err != nil {
		return nil, fmt.Errorf("INVALID_TRANSACTION_TYPE: The transaction type ID '%s' is invalid. Please provide a valid transaction type", req.TransactionType)
	}
	transactionTypeName, err := s.TransactionTypeService.GetTransactionTypeNameByID(transactionTypeID)
	if err != nil {
		return nil, fmt.Errorf("INVALID_TRANSACTION_TYPE: The transaction type does not exist")
	}

	relativeRate, err := s.getRelativeExchangeRate(req.FromCurrency, req.ToCurrency)
	if err != nil {
		return nil, err
	}

	transactionCode, err := s.CodeGen.GenerateCode()
	if err != nil {
		return nil, fmt.Errorf("TRANSACTION_CODE_ERROR: Failed to generate transaction code")
	}

	convertedAmount := utils.RoundOperations(req.Amount * relativeRate)

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	transactionsColl := database.MongoClient.Database(mongoDBName).Collection("transactions")
	transaction := models.Transaction{
		TransactionCode:   transactionCode,
		FromCurrency:      req.FromCurrency,
		ToCurrency:        req.ToCurrency,
		Amount:            req.Amount,
		AmountConverted:   convertedAmount,
		ExchangeRate:      relativeRate,
		TransactionTypeID: transactionTypeID,
		CreatedAt:         time.Now(),
		UserID:            userID,
	}

	result, err := transactionsColl.InsertOne(context.Background(), transaction)
	if err != nil {
		return nil, fmt.Errorf("DATABASE_ERROR: Failed to save transaction")
	}

	return map[string]interface{}{
		"transactionId":   result.InsertedID.(primitive.ObjectID).Hex(),
		"transactionCode": transaction.TransactionCode,
		"fromCurrency":    transaction.FromCurrency,
		"toCurrency":      transaction.ToCurrency,
		"amount":          transaction.Amount,
		"amountConverted": convertedAmount,
		"exchangeRate":    transaction.ExchangeRate,
		"transactionType": transactionTypeName,
		"createdAt":       transaction.CreatedAt.Format(time.RFC3339),
		"userId":          transaction.UserID,
	}, nil
}

func (s *ConversionService) isCurrencySupported(currency string) bool {
	for _, c := range s.SupportedCurrencies {
		if c == currency {
			return true
		}
	}
	return false
}

func (s *ConversionService) getRelativeExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", fromCurrency, toCurrency)
	relativeRate, err := database.RedisClient.Get(ctx, key).Float64()

	if err == redis.Nil {
		return 0, fmt.Errorf("EXCHANGE_RATE_ERROR: No  rate found for %s to %s", fromCurrency, toCurrency)
	} else if err != nil {
		return 0, fmt.Errorf("REDIS_ERROR: %v", err)
	}

	return relativeRate, nil
}
