package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunApiEngine() {
	engine := gin.Default()

	// OAuth2 Endpoints
	{
		router := engine.Group("/auth")
		{
			router := router.Group("/discord")
			router.GET("/login", handleLogin)
			router.GET("/callback", handleCallback)
		}
		router.POST("/logout", handleLogout)
	}

	// Api Endpoints
	router := engine.Group("/api")

	// Require authentication paths
	router.Use(authMiddleware())
	{
		router.GET("/me", func(ctx *gin.Context) { ctx.String(http.StatusOK, "Hello, you!") })
	}

	engine.Run()
}
