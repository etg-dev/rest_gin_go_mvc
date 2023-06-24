package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Posts    primitive.ObjectID `bson:"post,omitempty"`
	Action   []string           `bson:"action,omitempty"`
	ActionId primitive.ObjectID `bson:"actionId,omitempty"`
}
