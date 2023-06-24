package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/etg-dev/restApi/configs"
	"github.com/etg-dev/restApi/models"
	"github.com/etg-dev/restApi/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var actionCollection *mongo.Collection = configs.GetCollection(configs.DB, "actions")

// @descibe       Create new user
// @route         POST /users
// @access        Public
func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		newAction := models.Action{
			Actions: user.Action,
		}

		newUser := models.User{
			Name:  user.Name,
			Email: user.Email,
		}

		resultInsertedUser, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {

			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		userID, ok := resultInsertedUser.InsertedID.(primitive.ObjectID)
		if !ok {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "invalid ObjectID"}})
			return
		}
		newAction.User = userID

		resultInsertedAction, err := actionCollection.InsertOne(ctx, newAction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		actionID, ok := resultInsertedAction.InsertedID.(primitive.ObjectID)
		if !ok {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "invalid ObjectID"}})
			return
		}

		newUser.ActionId = actionID

		filter := bson.M{"_id": userID}
		update := bson.M{"$set": bson.M{"actionId": actionID}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedUser models.User
		err = userCollection.FindOne(ctx, filter).Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		responseUser := responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data: map[string]interface{}{
				"id":       updatedUser.Id,
				"name":     updatedUser.Name,
				"email":    updatedUser.Email,
				"actionId": updatedUser.ActionId,
			},
		}

		c.JSON(http.StatusCreated, responseUser)
	}
}

// @descibe       Get all users
// @route         GET /users
// @access        Public
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{}
		cur, err := userCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		for cur.Next(ctx) {
			var user models.User
			err := cur.Decode(&user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": users}})
	}
}

// @descibe       Get single user
// @route         GET /user/:id
// @access        Public
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": user}})
	}
}

// @descibe       Update single user
// @route         PUT /user/:id
// @access        Public
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		err = c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		filter := bson.M{"_id": id}
		updateField := bson.M{"$set": bson.M{"name": user.Name, "email": user.Email}}

		//! to return and update at the same time
		options := options.FindOneAndUpdate().SetReturnDocument(options.After)
		//
		err = userCollection.FindOneAndUpdate(ctx, filter, updateField, options).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": user}})
	}
}

// @descibe       Update single user
// @route         Delete /user/:id
// @access        Public
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Id not found", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		err = c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		filter := bson.M{"_id": id}

		result, err := userCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "User not found", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}
