package auth

import (
	"net/http"
	"quicksend/internal/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authService *Service, cfg *config.Config) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", loginHandler(authService))
		auth.POST("/callback", callbackHandler(authService))
		auth.POST("/refresh", RequireRefreshToken(authService), refreshTokenHandler(authService, cfg))
		auth.POST("/logout", RequireAccessToken(authService), logoutHandler)
	}
}

func loginHandler(authService *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authService.Login(c)
	}
}

func callbackHandler(authService *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authService.Login(c)
	}
}

func refreshTokenHandler(authService *Service, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		source := c.Query("source")

		refreshToken, _ := c.Get("refresh_token")

		tokenPair, err := authService.RefreshToken(refreshToken.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid refresh token: " + err.Error()})
			return
		}

		newAccessToken, newRefreshToken := tokenPair.AccessToken, tokenPair.RefreshToken

		if source == "website" {
			c.SetCookie("access_jwt_token", "Bearer "+newAccessToken,
				cfg.JWTAccessExpHours*3600,
				"/", "", true, true,
			)
			c.SetCookie("refresh_jwt_token", "Bearer "+newRefreshToken,
				cfg.JWTRefreshExpDays*3600*24,
				"/", "", true, true,
			)
			c.Status(http.StatusOK)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"access_jwt_token":  newAccessToken,
				"refresh_jwt_token": newRefreshToken,
			})
		}
	}
}

func logoutHandler(c *gin.Context) {
	c.SetCookie("access_jwt_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_jwt_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
