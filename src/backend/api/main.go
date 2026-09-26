package api

import (
	"net/http"

	"github.com/YutoMaeda1209/hygge/config"
	"github.com/gin-gonic/gin"
)

func RunApiEngine() error {
	engine := gin.Default()
	if err := engine.SetTrustedProxies(config.Conf.TrustProxyIp); err != nil {
		return err
	}

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
		router := router.Group("/user")
		router.GET("/me", func(ctx *gin.Context) { ctx.String(http.StatusOK, "Hello, you!") })
	}

	return engine.Run()
}
