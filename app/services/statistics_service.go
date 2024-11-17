package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatisticsService struct {
	DB *mongo.Collection
}

func NewStatisticsService(db *mongo.Database) *StatisticsService {
	return &StatisticsService{
		DB: db.Collection("transactions"),
	}
}

func (s *StatisticsService) GetStatisticsService(startDate, endDate time.Time) (map[string]interface{}, error) {

	pipeline := mongo.Pipeline{
		{{
			"$match", bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
			},
		}},
		{{
			"$lookup", bson.M{
				"from":         "transaction_types",
				"localField":   "transaction_type_id",
				"foreignField": "_id",
				"as":           "type_info",
			},
		}},
		{{
			"$unwind", "$type_info",
		}},
		{{
			"$group", bson.M{
				"_id": bson.M{
					"type_name": "$type_info.name",
					"currency":  "$to_currency",
				},
				"count":          bson.M{"$sum": 1},
				"totalConverted": bson.M{"$sum": "$amount_converted"},
				"totalAmount":    bson.M{"$sum": "$amount"},
			},
		}},
	}

	cursor, err := s.DB.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	response := map[string]interface{}{
		"totalTransactions":              0,
		"transactionsByType":             make(map[string]int),
		"totalAmountConvertedByCurrency": make(map[string]float64),
		"totalAmountByTransactionType":   make(map[string]float64),
		"averageAmountByTransactionType": make(map[string]float64),
	}

	countByType := make(map[string]int)
	totalAmountByType := make(map[string]float64)

	// Process results from cursor
	for cursor.Next(context.Background()) {
		var result struct {
			ID struct {
				TypeName string `bson:"type_name"`
				Currency string `bson:"currency"`
			} `bson:"_id"`
			Count          int     `bson:"count"`
			TotalConverted float64 `bson:"totalConverted"`
			TotalAmount    float64 `bson:"totalAmount"`
		}

		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		// Update response
		response["totalTransactions"] = response["totalTransactions"].(int) + result.Count
		response["transactionsByType"].(map[string]int)[result.ID.TypeName] += result.Count
		response["totalAmountConvertedByCurrency"].(map[string]float64)[result.ID.Currency] += result.TotalConverted
		response["totalAmountByTransactionType"].(map[string]float64)[result.ID.TypeName] += result.TotalAmount

		// save for average calculation
		countByType[result.ID.TypeName] += result.Count
		totalAmountByType[result.ID.TypeName] += result.TotalAmount
	}

	// Calculate averages
	for typeName, totalAmount := range totalAmountByType {
		count := countByType[typeName]
		if count > 0 {
			response["averageAmountByTransactionType"].(map[string]float64)[typeName] = totalAmount / float64(count)
		}
	}

	return response, nil
}
