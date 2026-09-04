package api

import (
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

	engine.Run()
}
