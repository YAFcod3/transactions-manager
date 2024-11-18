package services

import (
	"context"
	"time"
	"transactions-manager/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionsHistoryService struct {
	DB *mongo.Collection
}

func NewTransactionsHistoryService(db *mongo.Database) *TransactionsHistoryService {
	return &TransactionsHistoryService{
		DB: db.Collection("transactions"),
	}
}

func (s *TransactionsHistoryService) GetTransactionsHistory(ctx context.Context, startDate, endDate time.Time, transactionType string) (models.TransactionsHistoryResponse, error) {
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$lookup", bson.M{
			"from":         "transaction_types",
			"localField":   "transaction_type_id",
			"foreignField": "_id",
			"as":           "type_info",
		}}},
		{{"$unwind", bson.M{"path": "$type_info", "preserveNullAndEmptyArrays": true}}},
	}

	if transactionType != "" {
		pipeline = append(pipeline, bson.D{{"$match", bson.M{"type_info.name": transactionType}}})
	}

	pipeline = append(pipeline, bson.D{{"$project", bson.M{
		"transaction_code":    1,
		"from_currency":       1,
		"to_currency":         1,
		"amount":              1,
		"amount_converted":    1,
		"exchange_rate":       1,
		"transaction_type_id": 1,
		"created_at":          1,
		"user_id":             1,
		"transactionTypeName": "$type_info.name",
	}}})

	cursor, err := s.DB.Aggregate(ctx, pipeline)
	if err != nil {
		return models.TransactionsHistoryResponse{}, err
	}
	defer cursor.Close(ctx)

	var transactions []models.TransactionWithTypeInfo
	for cursor.Next(ctx) {
		var transaction models.TransactionWithTypeInfo
		if err := cursor.Decode(&transaction); err != nil {
			return models.TransactionsHistoryResponse{}, err
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return models.TransactionsHistoryResponse{}, err
	}

	response := models.TransactionsHistoryResponse{
		Total: len(transactions),
		Data:  []models.TransactionResponse{},
	}
	for _, transaction := range transactions {
		response.Data = append(response.Data, models.TransactionResponse{
			TransactionId:   transaction.TransactionID.Hex(),
			TransactionCode: transaction.TransactionCode,
			FromCurrency:    transaction.FromCurrency,
			ToCurrency:      transaction.ToCurrency,
			Amount:          transaction.Amount,
			AmountConverted: transaction.AmountConverted,
			ExchangeRate:    transaction.ExchangeRate,
			TransactionType: transaction.TransactionTypeName,
			CreatedAt:       transaction.CreatedAt.Format(time.RFC3339),
			UserId:          transaction.UserID,
		})
	}

	return response, nil
}
