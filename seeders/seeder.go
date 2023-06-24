package seeders

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/etg-dev/restApi/configs"
	"github.com/etg-dev/restApi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var postCollection *mongo.Collection = configs.GetCollection(configs.DB, "posts")

// SeedUsers inserts user data into the given collection
func InjectDB() error {
	var users []models.User
	var posts []models.Post

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	usersFilePath := path.Join(dir, "seeders", "users.json")
	postsFilePath := path.Join(dir, "seeders", "posts.json")

	//! Read data from user.json
	userData, err := os.ReadFile(usersFilePath)
	if err != nil {
		return err
	}

	//! Read data from post.json
	postData, err := os.ReadFile(postsFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(userData, &users)
	if err != nil {
		return err
	}

	var userDocs []interface{}
	for _, user := range users {
		user.Id = primitive.NewObjectID()
		userDocs = append(userDocs, user)
	}

	//! Inset data for users collection
	// _, err = userCollection.InsertMany(context.Background(), userDocs)
	// if err != nil {
	// 	return err
	// }
	//!

	err = json.Unmarshal(postData, &posts)
	if err != nil {
		return err
	}

	var postDocs []interface{}
	for _, post := range posts {
		post.User = primitive.NewObjectID()
		postDocs = append(postDocs, post)
	}

	_, err = postCollection.InsertMany(context.Background(), postDocs)
	if err != nil {
		return err
	}

	fmt.Println("DATA ADDED...")
	return nil
}

// DeleteUsers deletes all data from the given collection
func DrainDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	_, err = postCollection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	fmt.Println("DATA DELETED...")
	return nil
}
