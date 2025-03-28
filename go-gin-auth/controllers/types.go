package controllers

import (
	"github.com/CompileWithG/go-gin-auth/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// HashResult holds the result of a password hashing operation
type HashResult struct {
	Hashed string
	Err    error
}

// InsertResult holds the result of a database insert operation
type InsertResult struct {
	Result *mongo.InsertOneResult
	Err    error
}

// UserResult holds the result of a user lookup operation
type UserResult struct {
	User models.User
	Err  error
}

// TokenResult holds the result of a token generation operation
type TokenResult struct {
	Token string
	Err   error
}
