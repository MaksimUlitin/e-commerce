package middleware

import (
	"github.com/gin-gonic/gin"
	jwt2 "github.com/maksimulitin/internal/jwt"
	"net/http"
)

func middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")

		if clientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
			ctx.Abort()
			return
		}
		claims, err := jwt2.ValidateToken(clientToken)
		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			ctx.Abort()
		}
		ctx.Set("email", claims.Email)
		ctx.Set("claims", claims.Uid)
	}
}
