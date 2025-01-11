package middleware

import (
	token "github.com/maksimulitin/internal/tokens"
	"github.com/maksimulitin/lib/logger"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			logger.Error("no token provided", slog.Any("token", c.Request.Header.Get("token")))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			logger.Error("invalid token", slog.Any("token", c.Request.Header.Get("token")))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
