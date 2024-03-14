package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func login(c *gin.Context) {
	_, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not found"})
	}
	// Do something with the database

	// Return a response
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
