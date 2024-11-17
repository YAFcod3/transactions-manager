package services

import (
	"context"
	"errors"
	"os"
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
}

func NewConversionService(codeGen *generate_transaction_code.CodeGenerator, supportedCurrenciesService *SupportedCurrenciesService, transactionTypeService *TransactionTypeService) *ConversionService {
	return &ConversionService{
		CodeGen:                    codeGen,
		SupportedCurrenciesService: supportedCurrenciesService,
		TransactionTypeService:     transactionTypeService,
	}
}

func (s *ConversionService) ProcessTransaction(req models.TransactionRequest, userID string) (map[string]interface{}, error) {
	mongoDBName := os.Getenv("MONGO_DB_NAME")

	if err := s.SupportedCurrenciesService.IsCurrencySupported(req.FromCurrency); err != nil {
		return nil, errors.New("the 'fromCurrency' is not supported")
	}
	if err := s.SupportedCurrenciesService.IsCurrencySupported(req.ToCurrency); err != nil {
		return nil, errors.New("the 'toCurrency' is not supported")
	}

	transactionTypeID, err := primitive.ObjectIDFromHex(req.TransactionType)
	if err != nil {
		return nil, errors.New("invalid transaction type")
	}

	transactionTypeName, err := s.TransactionTypeService.GetTransactionTypeNameByID(transactionTypeID)
	if err != nil {
		return nil, err
	}

	fromRate, err := utils.GetExchangeRate(req.FromCurrency)
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, errors.New("'fromCurrency' exchange rate not found in Redis")
		}
		return nil, errors.New("failed to fetch 'fromCurrency' exchange rate")
	}

	toRate, err := utils.GetExchangeRate(req.ToCurrency)
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, errors.New("'toCurrency' exchange rate not found in Redis")
		}
		return nil, errors.New("failed to fetch 'toCurrency' exchange rate")
	}

	transactionCode, err := s.CodeGen.GenerateCode()
	if err != nil {
		return nil, errors.New("failed to generate transaction code")
	}

	convertedAmount := (req.Amount / fromRate) * toRate

	transaction := models.Transaction{
		TransactionCode:   transactionCode,
		FromCurrency:      req.FromCurrency,
		ToCurrency:        req.ToCurrency,
		Amount:            req.Amount,
		AmountConverted:   convertedAmount,
		ExchangeRate:      toRate / fromRate,
		TransactionTypeID: transactionTypeID,
		CreatedAt:         time.Now(),
		UserID:            userID,
	}

	transactionsColl := database.MongoClient.Database(mongoDBName).Collection("transactions")
	result, err := transactionsColl.InsertOne(context.Background(), transaction)
	if err != nil {
		return nil, errors.New("failed to save transaction")
	}

	return map[string]interface{}{
		"transactionId":   result.InsertedID.(primitive.ObjectID).Hex(),
		"transactionCode": transaction.TransactionCode,
		"fromCurrency":    transaction.FromCurrency,
		"toCurrency":      transaction.ToCurrency,
		"amount":          transaction.Amount,
		"amountConverted": transaction.AmountConverted,
		"exchangeRate":    transaction.ExchangeRate,
		"transactionType": transactionTypeName,
		"createdAt":       transaction.CreatedAt.Format(time.RFC3339),
		"userId":          transaction.UserID,
	}, nil
}
