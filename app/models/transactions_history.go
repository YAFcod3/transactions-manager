package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionResponse struct {
	TransactionId   string  `json:"transactionId"`
	TransactionCode string  `json:"transactionCode"`
	FromCurrency    string  `json:"fromCurrency"`
	ToCurrency      string  `json:"toCurrency"`
	Amount          float64 `json:"amount"`
	AmountConverted float64 `json:"amountConverted"`
	ExchangeRate    float64 `json:"exchangeRate"`
	TransactionType string  `json:"transactionType"`
	CreatedAt       string  `json:"createdAt"`
	UserId          string  `json:"userId"`
}

type TransactionWithTypeInfo struct {
	TransactionID       primitive.ObjectID `bson:"_id"`
	TransactionCode     string             `bson:"transaction_code"`
	FromCurrency        string             `bson:"from_currency"`
	ToCurrency          string             `bson:"to_currency"`
	Amount              float64            `bson:"amount"`
	AmountConverted     float64            `bson:"amount_converted"`
	ExchangeRate        float64            `bson:"exchange_rate"`
	TransactionTypeID   primitive.ObjectID `bson:"transaction_type_id"`
	TransactionTypeName string             `bson:"transactionTypeName"`
	CreatedAt           time.Time          `bson:"created_at"`
	UserID              string             `bson:"user_id"`
}

type TransactionsHistoryResponse struct {
	Total int                   `json:"total"`
	Data  []TransactionResponse `json:"data"`
}
