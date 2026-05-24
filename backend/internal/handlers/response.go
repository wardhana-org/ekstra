package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeInvalidRequest(c *gin.Context) {
	writeError(c, http.StatusBadRequest, "invalid request")
}

func writeUnauthenticated(c *gin.Context) {
	writeError(c, http.StatusUnauthorized, "unauthenticated")
}

func writeInternalServerError(c *gin.Context) {
	writeError(c, http.StatusInternalServerError, "internal server error")
}

func writeError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}
