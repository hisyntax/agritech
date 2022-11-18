package persistence

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	unAuth := errors.New("unauthorized").Error()
	clientToken := c.Request.Header.Get("Authorization")
	if clientToken == "" {
		msg := "No Authorization header provided"
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		c.Abort()
		return
	}

	claims, msg := ValidateToken(clientToken)
	if msg != "" {
		c.JSON(http.StatusBadRequest, unAuth)
		c.Abort()
		return
	}

	c.Set("uid", claims.Uid)
	c.Next()
}
