package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type TransactionType struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Name        string             `bson:"name" json:"name"`
// 	Description string             `bson:"description,omitempty" json:"description"`
// }

type TransactionType struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Deleted     bool               `bson:"deleted" json:"-"`
}
