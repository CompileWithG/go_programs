package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// API routes
	api := router.Group("/api")
	{
		api.GET("/health", HealthCheck())
		api.POST("/signup", CreateUser())
		api.POST("/login", GetAUser())
	}

	return router
}
