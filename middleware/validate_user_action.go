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
	"go.mongodb.org/mongo-driver/mongo"
)

var actionCollection *mongo.Collection = configs.GetCollection(configs.DB, "actions")

func ValidateAction(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Get userId from context
		userId, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Invalid userId"})
			c.Abort()
			return
		}

		// Find user's actions
		filterActions := bson.M{"user": userId}
		var foundUserAction models.Action
		err := actionCollection.FindOne(ctx, filterActions).Decode(&foundUserAction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "There is no action for that user", Data: map[string]interface{}{"data": err.Error()}})
			c.Abort()
			return
		}

		// Check if user has the required action
		hasAction := false
		for _, a := range foundUserAction.Actions {
			if a == action {
				hasAction = true
				break
			}
		}

		if !hasAction {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: fmt.Sprintf("This user does not have access to %s", action)})
			c.Abort()
			return
		}

		c.Next()
	}
}
