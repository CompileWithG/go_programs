package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck endpoint
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "API is running",
			"data":    nil,
		})

	}
}
