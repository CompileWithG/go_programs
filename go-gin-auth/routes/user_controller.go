package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/CompileWithG/go-gin-auth/config"
	"github.com/CompileWithG/go-gin-auth/controllers"
	"github.com/CompileWithG/go-gin-auth/models"
	"github.com/CompileWithG/go-gin-auth/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}
		hashed_password, err1 := controllers.HashPassword(user.Password)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error hashing password", Data: map[string]interface{}{"data": err1.Error()}})
			return
		}
		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Name:     user.Name,
			Email:    user.Email,
			Password: hashed_password,
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}
func GetAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var loginRequest models.LoginRequest
		var user models.User
		defer cancel()

		// Validate request body
		if err := c.BindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		// Find user by email
		err := userCollection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": "Invalid email or password"},
			})
			return
		}

		// Verify password (assuming you have a password hashing mechanism)
		// You should have a function to compare hashed password with input password
		if !controllers.CheckPassword(loginRequest.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, responses.UserResponse{
				Status:  http.StatusUnauthorized,
				Message: "error",
				Data:    map[string]interface{}{"data": "Invalid email or password"},
			})
			return
		}

		// Generate JWT token
		token, err := controllers.GenerateToken(user.Id.Hex())

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": "Failed to generate token"},
			})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"token": token, "user": user},
		})
	}
}
