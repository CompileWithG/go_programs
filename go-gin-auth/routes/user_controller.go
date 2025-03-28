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
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		// Validate user struct in a goroutine
		validationErrChan := make(chan error, 1)
		go func() {
			validationErrChan <- validate.Struct(&user)
		}()

		// Hash password in parallel with validation
		hashChan := make(chan controllers.HashResult, 1)
		go func() {
			hashed, err := controllers.HashPassword(user.Password)
			hashChan <- controllers.HashResult{Hashed: hashed, Err: err}
		}()

		// Wait for validation
		if validationErr := <-validationErrChan; validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": validationErr.Error()},
			})
			return
		}

		// Wait for password hash
		hashResult := <-hashChan
		if hashResult.Err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error hashing password",
				Data:    map[string]interface{}{"data": hashResult.Err.Error()},
			})
			return
		}

		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Name:     user.Name,
			Email:    user.Email,
			Password: hashResult.Hashed,
		}

		// Insert user in a goroutine
		insertChan := make(chan controllers.InsertResult, 1)
		go func() {
			result, err := userCollection.InsertOne(ctx, newUser)
			insertChan <- controllers.InsertResult{Result: result, Err: err}
		}()

		insertResult := <-insertChan
		if insertResult.Err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": insertResult.Err.Error()},
			})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data:    map[string]interface{}{"data": insertResult.Result},
		})
	}
}

func GetAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var loginRequest models.LoginRequest
		if err := c.BindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		// Find user by email in a goroutine
		userChan := make(chan controllers.UserResult, 1)
		go func() {
			var user models.User
			err := userCollection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&user)
			userChan <- controllers.UserResult{User: user, Err: err}
		}()

		userResult := <-userChan
		if userResult.Err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": "Invalid email or password"},
			})
			return
		}

		// Verify password in a goroutine
		passwordChan := make(chan bool, 1)
		go func() {
			passwordChan <- controllers.CheckPassword(loginRequest.Password, userResult.User.Password)
		}()

		if !<-passwordChan {
			c.JSON(http.StatusUnauthorized, responses.UserResponse{
				Status:  http.StatusUnauthorized,
				Message: "error",
				Data:    map[string]interface{}{"data": "Invalid email or password"},
			})
			return
		}

		// Generate token in a goroutine
		tokenChan := make(chan controllers.TokenResult, 1)
		go func() {
			token, err := controllers.GenerateToken(userResult.User.Id.Hex())
			tokenChan <- controllers.TokenResult{Token: token, Err: err}
		}()

		tokenResult := <-tokenChan
		if tokenResult.Err != nil {
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
			Data:    map[string]interface{}{"token": tokenResult.Token, "user": userResult.User},
		})
	}
}
