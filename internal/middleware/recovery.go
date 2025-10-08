package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery recovers from panics and returns a 500 error
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

// ErrorHandler handles errors and returns consistent error responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Return error response
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal server error",
				"details": err.Error(),
			})
		}
	}
}
