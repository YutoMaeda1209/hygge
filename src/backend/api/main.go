package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Api() {
	auth()

	engine := gin.Default()

	// Api Endpoints
	router := engine.Group("/api")
	{
		router := router.Group("/auth")
		router.GET("/login", handleLogin)
		router.GET("/callback", handleCallback)
	}

	// Require authentication paths
	router.Use(AuthMiddleware())
	{
		router := router.Group("/hello")
		router.GET("/world", func(ctx *gin.Context) { ctx.String(http.StatusOK, "Hello, world!") })
		router.GET("/you", func(ctx *gin.Context) { ctx.String(http.StatusOK, "Hello, you!") })
	}

	engine.Run()
}
