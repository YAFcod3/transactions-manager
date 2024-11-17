package services

import (
	"context"
	"errors"
	"transactions-manager/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionTypeService struct {
	DB *mongo.Collection
}

func NewTransactionTypeService(db *mongo.Database) *TransactionTypeService {
	return &TransactionTypeService{
		DB: db.Collection("transaction_types"),
	}
}

func (s *TransactionTypeService) CreateTransactionType(transactionType models.TransactionType) error {
	count, err := s.DB.CountDocuments(context.Background(), bson.M{
		"name": bson.M{"$regex": transactionType.Name, "$options": "i"},
	})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("transaction type already exists")
	}

	_, err = s.DB.InsertOne(context.Background(), transactionType)
	return err
}

func (s *TransactionTypeService) GetTransactionTypes() ([]models.TransactionType, error) {
	cursor, err := s.DB.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var types []models.TransactionType
	if err := cursor.All(context.Background(), &types); err != nil {
		return nil, err
	}
	return types, nil
}

func (s *TransactionTypeService) GetTransactionTypeByID(id string) (models.TransactionType, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.TransactionType{}, errors.New("invalid ID")
	}
	var transactionType models.TransactionType
	err = s.DB.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&transactionType)
	if err == mongo.ErrNoDocuments {
		return models.TransactionType{}, errors.New("transaction type not found")
	}
	return transactionType, err
}

func (s *TransactionTypeService) UpdateTransactionType(id string, updatedType models.TransactionType) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID")
	}

	if updatedType.Name != "" {
		count, err := s.DB.CountDocuments(context.Background(), bson.M{
			"name": bson.M{"$regex": updatedType.Name, "$options": "i"},
			"_id":  bson.M{"$ne": objectID},
		})
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("transaction type name already exists")
		}
	}

	updateFields := bson.M{}
	if updatedType.Name != "" {
		updateFields["name"] = updatedType.Name
	}
	if updatedType.Description != "" {
		updateFields["description"] = updatedType.Description
	}

	_, err = s.DB.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updateFields})
	return err
}

func (s *TransactionTypeService) DeleteTransactionType(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID")
	}
	_, err = s.DB.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

func (s *TransactionTypeService) GetTransactionTypeNameByID(id primitive.ObjectID) (string, error) {
	var transactionType models.TransactionType
	err := s.DB.FindOne(context.Background(), bson.M{"_id": id}).Decode(&transactionType)
	if err == mongo.ErrNoDocuments {
		return "", errors.New("transaction type not found")
	}
	return transactionType.Name, err
}
