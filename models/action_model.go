package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Action struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	User    primitive.ObjectID `bson:"user,omitempty"`
	Actions []string           `bson:"actions,omitempty"`
}
