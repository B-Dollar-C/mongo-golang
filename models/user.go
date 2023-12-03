package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	AuthToken string             `json: "authToken,omitempty, bson: "authToken,omitempty`
	Name      string             `json: "name" bson: "name"`
	Gender    string             `json: "gender" bson: "gender"`
	Age       int                `json: "age" bson: "age"`
	Id        primitive.ObjectID `json: "id,omitempty" bson: "_id,omitempty"`
}

type InsertedId struct {
	InsertedId string
}
