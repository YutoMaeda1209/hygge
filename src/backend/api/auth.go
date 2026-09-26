package api

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/YutoMaeda1209/hygge/config"
	"github.com/YutoMaeda1209/hygge/model"
	"github.com/gin-gonic/gin"
)

const stateCookieName = "oauth_state"
const stateCookieAge = time.Minute * 5 // 5 minutes

const errMsgRetry = "エラーが発生しました。もう一度やり直してください。"
const errMsgRetryLater = "エラーが発生しました。時間を置いてやり直してください。"

func handleLogin(c *gin.Context) {
	// Generate and store state data
	state, err := generateState()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate state")
		slog.Error("Failed to generate state", "err", err)
		return
	}
	setCookie(c, stateCookieName, state, int(stateCookieAge/time.Second))

	url := config.Conf.OAuth2Config.AuthCodeURL(state)
	c.Redirect(302, url)
}

func handleCallback(c *gin.Context) {
	// Verify and discard state data
	state, err := c.Cookie(stateCookieName)
	if err != nil || state == "" || c.Query("state") != state {
		redirectErrPage(c, errMsgRetry)
		slog.Warn("Invalid state", "err", err)
		return
	}
	setCookie(c, stateCookieName, "", -1)

	// Verify oauth2 authorization code
	code := c.Query("code")
	if code == "" {
		redirectErrPage(c, errMsgRetry)
		slog.Warn("Missing code")
		return
	}

	// Get token from discord
	token, err := config.Conf.OAuth2Config.Exchange(c.Request.Context(), code)
	if err != nil {
		redirectErrPage(c, errMsgRetry)
		slog.Warn("Failed to exchange token", "err", err)
		return
	}

	// Fetch the discord user tied to the token
	client := config.Conf.OAuth2Config.Client(c.Request.Context(), token)
	identify, err := fetchDiscordIdentify(client)
	if err != nil {
		redirectErrPage(c, errMsgRetryLater)
		slog.Error("Failed to fetch discord user", "err", err)
		return
	}

	// Create or obtain an account
	account := model.Account{}
	if err := account.EnsureByDiscordId(c.Request.Context(), identify.Id); err != nil {
		redirectErrPage(c, errMsgRetry)
		slog.Error("Can not ensure an account", "err", err)
		return
	}

	// Start a session for the user
	if err := startSession(c, account.Id); err != nil {
		redirectErrPage(c, errMsgRetry)
		slog.Error("Failed to start a session", "err", err)
		return
	}

	c.Redirect(302, "/dashboard")
}

func handleLogout(c *gin.Context) {
	if token, err := c.Cookie(sessionCookieName); err == nil && token != "" {
		if err := model.DeleteSession(c.Request.Context(), token); err != nil {
			c.String(http.StatusInternalServerError, "Failed to logout")
			slog.Error("Failed to delete session", "err", err)
			return
		}
	}
	removeCookie(c, sessionCookieName)
	c.Redirect(302, "/")
}

func generateState() (string, error) {
	const stateLen = 32
	b := make([]byte, stateLen/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func redirectErrPage(c *gin.Context, msg string) {
	url := url.URL{Path: "/signin"}
	query := url.Query()
	query.Set("err", msg)
	url.RawQuery = query.Encode()
	c.Redirect(302, url.String())
}
