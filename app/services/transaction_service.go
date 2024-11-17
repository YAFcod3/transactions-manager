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
	CodeGen *generate_transaction_code.CodeGenerator
}

func NewConversionService(codeGen *generate_transaction_code.CodeGenerator) *ConversionService {
	return &ConversionService{
		CodeGen: codeGen,
	}
}

func (s *ConversionService) ProcessTransaction(req models.TransactionRequest, userID string) (map[string]interface{}, error) {
	mongoDBName := os.Getenv("MONGO_DB_NAME")

	if req.Amount <= 0 || req.FromCurrency == "" || req.ToCurrency == "" || req.TransactionType == "" {
		return nil, errors.New("invalid input")
	}

	// Validar el tipo de transacción en MongoDB
	// transactionTypeID, err := primitive.ObjectIDFromHex(req.TransactionType)
	// if err != nil {
	// 	return nil, errors.New("invalid transaction type")
	// }
	// collection :=  database.MongoClient.Database(mongoDBName).Collection("transaction_types")
	// var transactionType struct {
	// 	Name string `bson:"name"`
	// }
	// err = collection.FindOne(context.Background(), bson.M{"_id": transactionTypeID}).Decode(&transactionType)
	// if err != nil {
	// 	return nil, errors.New("transaction type not found")
	// }

	// get rate from redis
	fromRate, err := utils.GetExchangeRate(req.FromCurrency)
	if err != nil {
		return nil, errors.New("from currency not found")
	}
	// toRate, err := s.getExchangeRate(req.ToCurrency)
	toRate, err := utils.GetExchangeRate(req.ToCurrency)
	if err != nil {
		return nil, errors.New("to currency not found")
	}

	transactionCode, err := s.CodeGen.GenerateCode()
	if err != nil {
		return nil, errors.New("failed to generate transaction code")
	}

	convertedAmount := (req.Amount / fromRate) * toRate

	transaction := models.Transaction{
		TransactionCode: transactionCode,
		FromCurrency:    req.FromCurrency,
		ToCurrency:      req.ToCurrency,
		Amount:          req.Amount,
		AmountConverted: convertedAmount,
		ExchangeRate:    toRate / fromRate,
		// TransactionTypeID: transactionTypeID,
		CreatedAt: time.Now(),
		UserID:    userID,
	}

	// transactionsColl := s.MongoClient.Database("currencyMongoDb").Collection("transactions")
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
		// "transactionType": transactionType.Name,
		"createdAt": transaction.CreatedAt.Format(time.RFC3339),
		"userId":    transaction.UserID,
	}, nil
}
