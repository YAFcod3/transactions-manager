package utils

import (
	"context"
	"errors"
	"strconv"
	"transactions-manager/app/database"
)

func GetExchangeRate(currency string) (float64, error) {
	if currency == "USD" {
		return 1.0, nil
	}
	ctx := context.Background()
	rateStr, err := database.RedisClient.HGet(ctx, "exchange_rates", currency).Result()
	if err != nil {
		return 0, errors.New("currency not found")
	}
	return strconv.ParseFloat(rateStr, 64)
}
