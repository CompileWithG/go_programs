package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/CompileWithG/go-gin-auth/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	
	// API routes
	api := router.Group("/api")
	{
		api.GET("/health", controllers.HealthCheck())
		api.POST("/signup", controllers.CreateUser())
    api.POST("/login",controllers.GetAUser())
	}

	return router
}
