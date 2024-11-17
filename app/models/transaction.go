package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionRequest struct {
	FromCurrency    string  `json:"fromCurrency"`
	ToCurrency      string  `json:"toCurrency"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transactionType"`
}

type Transaction struct {
	TransactionCode   string             `bson:"transaction_code"`
	FromCurrency      string             `bson:"from_currency"`
	ToCurrency        string             `bson:"to_currency"`
	Amount            float64            `bson:"amount"`
	AmountConverted   float64            `bson:"amount_converted"`
	ExchangeRate      float64            `bson:"exchange_rate"`
	TransactionTypeID primitive.ObjectID `bson:"transaction_type_id"`
	CreatedAt         time.Time          `bson:"created_at"`
	UserID            string             `bson:"user_id"`
}
