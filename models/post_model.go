package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `bson:"title,omitempty"`
	Content string             `bson:"content,omitempty"`
	User    primitive.ObjectID `bson:"user,omitempty"`
}
