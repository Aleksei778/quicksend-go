package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireRefreshToken(authService *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rawToken string

		cookie, err := c.Cookie("refresh_jwt_token")
		if err == nil {
			rawToken = cookie
		} else {
			rawToken = c.GetHeader("Authorization")
		}

		tokenStr := strings.TrimPrefix(rawToken, "Bearer ")

		claims, err := authService.VerifyRefreshToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid refresh token"})
			return
		}

		c.Set("refresh_token", claims)
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func RequireAccessToken(authService *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rawToken string
		cookie, err := c.Cookie("access_jwt_token")
		if err == nil {
			rawToken = cookie
		} else {
			rawToken = c.GetHeader("Authorization")
		}

		tokenStr := strings.TrimPrefix(rawToken, "Bearer ")

		claims, err := authService.VerifyAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid access token"})
			return
		}

		c.Set("access_token", claims)
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
