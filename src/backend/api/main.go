package api

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Api() {
	auth()

	// Setup db
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		slog.Error("The database URL has not been set. Visit https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING-KEYWORD-VALUE to configure it.")
		os.Exit(1)
	}

	if postgresDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		slog.Error("Failed to connect to the database", "err", err)
		os.Exit(1)
	} else {
		db = postgresDb
	}

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
