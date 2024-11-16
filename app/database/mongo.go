package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitMongo() {
	mongoUsername := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoDBHost := "transactionsManagerDb"
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	mongoPort := "27017"

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUsername, mongoPassword, mongoDBHost, mongoPort)
	log.Printf("Connecting to MongoDB at %s", mongoURI)

	var err error
	MongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := MongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Printf("Connected to MongoDB database: %s", mongoDBName)
}

func CloseMongo() {
	if err := MongoClient.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
	} else {
		log.Println("MongoDB connection closed")
	}
}
