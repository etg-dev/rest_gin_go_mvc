package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/etg-dev/restApi/configs"
	"github.com/etg-dev/restApi/models"
	"github.com/etg-dev/restApi/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func ValidateUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Get User ID from url
		userId, err := primitive.ObjectIDFromHex(c.Param("userId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Id required", Data: map[string]interface{}{"data": err.Error()}})
			c.Abort()
			return
		}

		// Check user Id exist
		filterUser := bson.M{"_id": userId}
		var foundUser models.User
		err = userCollection.FindOne(ctx, filterUser).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "User not found with that id", Data: map[string]interface{}{"data": err.Error()}})
			c.Abort()
			return
		}
		fmt.Println(filterUser)

		c.Set("userId", userId) // set userId in context
		c.Next()
	}
}
