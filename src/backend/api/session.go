package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/YutoMaeda1209/hygge/model"
	"github.com/gin-gonic/gin"
)

const sessionCookieName = "session"
const sessionAge = time.Hour * 24 * 7          // 7 days
const sessionRefreshThreshold = sessionAge / 2 // 3.5 days
const sessionCleanupInterval = time.Hour

func startSession(c *gin.Context, accountId uint) error {
	token, err := model.CreateSession(c.Request.Context(), accountId, sessionAge)
	if err != nil {
		return err
	}
	setCookie(c, sessionCookieName, token, int(sessionAge/time.Second))
	return nil
}

func cleanupExpiredSessions() {
	for {
		if n, err := model.DeleteExpiredSessions(context.Background()); err != nil {
			slog.Error("Failed to delete expired sessions", "err", err)
		} else if n > 0 {
			slog.Info("Deleted expired sessions", "count", n)
		}
		time.Sleep(sessionCleanupInterval)
	}
}
