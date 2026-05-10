package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success sends a 200 (or custom status) JSON response with the given data payload.
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// Created sends a 201 JSON response with the given data payload.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": data})
}

// NoContent sends a 204 response with no body.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends a JSON error response using the provided HTTP status code and message.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"success": false, "error": message})
}
