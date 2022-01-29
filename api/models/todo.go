package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Task        string             `json:"task" bson:"task"`
	IsCompleted bool               `json:"is_completed" bson:"is_completed"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
}
