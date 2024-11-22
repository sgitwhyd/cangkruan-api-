package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/cangkruan-api/internal/configs"
	"github.com/sgitwhyd/cangkruan-api/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	secretKey := configs.Get().Service.SecretJWT

	return func(ctx *gin.Context) {
		header :=  ctx.Request.Header.Get("Authorization")
		header = strings.TrimSpace(header)
		if header == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.New("invalid token"),
			})
			return
		}

		userID, username, err := jwt.ValidateToken(header, secretKey)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.New("invalid token"),
			})
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("username", username)
		ctx.Next()
	}
}