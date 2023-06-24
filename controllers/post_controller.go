package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/etg-dev/restApi/configs"
	"github.com/etg-dev/restApi/models"
	"github.com/etg-dev/restApi/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection = configs.GetCollection(configs.DB, "posts")

// @descibe       Create new post
// @route         POST /posts/:userId
// @access        Public
func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var input models.Post
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusInternalServerError, responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if input.Title == "" || input.Content == "" {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "title and content required"}})
			return
		}

		userId, err := primitive.ObjectIDFromHex(c.Param("userId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Id required", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Check if the user exists
		userFilter := bson.M{"_id": userId}
		user, err := postCollection.Find(ctx, userFilter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if user == nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusInternalServerError, Message: "User not found", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		post := models.Post{
			Title:   input.Title,
			Content: input.Content,
			User:    userId,
		}

		result, err := postCollection.InsertOne(ctx, post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

// @descibe       Get all posts
// @route         GET /posts/:userId
// @access        Public
// func GetPosts() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		var posts []models.Post
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()

// 		page, ok := c.Get("page")
// 		if !ok {
// 			page = 1
// 		}

// 		pageSize, ok := c.Get("pageSize")
// 		if !ok {
// 			pageSize = 10
// 		}

// 		filter := bson.M{}

// 		projection := bson.M{}

// 		fieldsParam := c.Query("select")
// 		if fieldsParam != "" {
// 			fields := strings.Split(fieldsParam, ",")
// 			for _, field := range fields {
// 				projection[field] = 1
// 				fmt.Println(projection[field])
// 			}
// 		}

// 		options := options.Find().SetLimit(pageSize.(int64)).SetSkip(page.(int64) * pageSize.(int64)).SetProjection(projection)

// 		cur, err := postCollection.Find(ctx, filter, options)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
// 			return
// 		}

// 		for cur.Next(ctx) {
// 			var post models.Post
// 			err := cur.Decode(&post)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}
// 			posts = append(posts, post)
// 		}
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
// 			return
// 		}

// 		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": posts}})
// 	}
// }

func GetPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		page, ok := c.Get("page")
		if !ok {
			page = 1
		}

		pageSize, ok := c.Get("pageSize")
		if !ok {
			pageSize = 10
		}

		fieldsParam := c.Query("select")
		projection := bson.M{}
		if fieldsParam != "" {
			fields := strings.Split(fieldsParam, ",")
			for _, field := range fields {
				projection[field] = 1
			}
		} else {
			projection = bson.M{"id": 0}
		}

		pipeline := []bson.M{
			{"$match": bson.M{}},
			{"$sort": bson.M{"created_at": -1}},
			{"$skip": page.(int64) * pageSize.(int64)},
			{"$limit": pageSize.(int64)},
			{"$project": projection},
		}

		cursor, err := postCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}
		defer cursor.Close(ctx)

		var posts []bson.M
		if err = cursor.All(ctx, &posts); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"data": posts},
		})
	}
}

// @descibe       Get all user posts
// @route         GET /posts/user/:userId
// @access        Public
func GetUsersPosts() gin.HandlerFunc {

	return func(c *gin.Context) {
		var posts []models.Post
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userId, err := primitive.ObjectIDFromHex(c.Param("userId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusInternalServerError, Message: "User id required", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		fmt.Println(userId)

		filter := bson.M{"user": userId}
		cur, err := postCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		for cur.Next(ctx) {
			var post models.Post
			err := cur.Decode(&post)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			posts = append(posts, post)
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": posts}})
	}
}
