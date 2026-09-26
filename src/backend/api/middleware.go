package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/YutoMaeda1209/hygge/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const accountContextKey = "account"

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(sessionCookieName)
		if err != nil || token == "" {
			c.String(http.StatusUnauthorized, "Missing session")
			c.Abort()
			return
		}

		session, err := model.FindSession(c.Request.Context(), token)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			removeCookie(c, sessionCookieName)
			c.String(http.StatusUnauthorized, "Invalid session")
			c.Abort()
			return
		} else if err != nil {
			slog.Error("Failed to find session", "err", err)
			c.String(http.StatusInternalServerError, "")
			c.Abort()
			return
		}

		// Sliding expiration; extending is idempotent, so concurrent requests are safe.
		if time.Until(session.ExpiresAt) < sessionRefreshThreshold {
			if err := session.Extend(c.Request.Context(), sessionAge); err != nil {
				slog.Error("Failed to extend session", "err", err)
			} else {
				setCookie(c, sessionCookieName, token, int(sessionAge/time.Second))
			}
		}

		c.Set(accountContextKey, session.Account)
		c.Next()
	}
}
