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

	fromRate, err := utils.GetExchangeRate(req.FromCurrency)
	if err != nil {
		return nil, fmt.Errorf("EXCHANGE_RATE_ERROR: 'fromCurrency' exchange rate not found")
	}
	toRate, err := utils.GetExchangeRate(req.ToCurrency)
	if err != nil {
		return nil, fmt.Errorf("EXCHANGE_RATE_ERROR: 'toCurrency' exchange rate not found")
	}

	transactionCode, err := s.CodeGen.GenerateCode()
	if err != nil {
		return nil, fmt.Errorf("TRANSACTION_CODE_ERROR: Failed to generate transaction code")
	}

	convertedAmount := (req.Amount / fromRate) * toRate

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	transactionsColl := database.MongoClient.Database(mongoDBName).Collection("transactions")
	transaction := models.Transaction{
		TransactionCode:   transactionCode,
		FromCurrency:      req.FromCurrency,
		ToCurrency:        req.ToCurrency,
		Amount:            req.Amount,
		AmountConverted:   utils.RoundOperations(convertedAmount),
		ExchangeRate:      utils.RoundOperations(toRate / fromRate),
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
		"amountConverted": utils.RoundOperations(transaction.AmountConverted),
		"exchangeRate":    utils.RoundOperations(transaction.ExchangeRate),
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
