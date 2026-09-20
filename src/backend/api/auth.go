package api

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/YutoMaeda1209/hygge/config"
	"github.com/YutoMaeda1209/hygge/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

const stateCookieName = "oauth_state"
const stateCookieAge = time.Minute * 5 // 5 minutes
const jwtCookieName = "jwt"
const jwtTokenAge = time.Hour * 24 * 7 // 7 days

type claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func handleLogin(c *gin.Context) {
	// Generate and store state data
	state, err := generateState()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate state")
		slog.Error("Failed to generate state", "err", err)
		return
	}
	setCookie(c, stateCookieName, state, int(stateCookieAge/time.Second))

	url := config.Conf.OAuth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(302, url)
}

func handleCallback(c *gin.Context) {
	// Verify and discard state data
	state, err := c.Cookie(stateCookieName)
	if err != nil || c.Query("state") != state {
		redirectErrPage(c)
		slog.Warn("Invalid state", "err", err)
		return
	}
	setCookie(c, stateCookieName, "", -1)

	// Verify oauth2 authorization code
	code := c.Query("code")
	if code == "" {
		redirectErrPage(c)
		slog.Warn("Missing code")
		return
	}

	// Get token from discord
	token, err := config.Conf.OAuth2Config.Exchange(c.Request.Context(), code)
	if err != nil {
		redirectErrPage(c)
		slog.Warn("Failed to exchange token", "err", err)
		return
	}

	// Fetch the discord user tied to the token
	client := config.Conf.OAuth2Config.Client(c.Request.Context(), token)
	discordIdentify, err := model.FetchDiscordIdentify(client)
	if err != nil {
		redirectErrPage(c, "エラーが発生しました。時間を置いてやり直してください。")
		slog.Error("Failed to fetch discord user", "err", err)
		return
	}

	// Create or obtain an account
	account := model.Account{}
	err = account.Find(c, "discord_id = ?", discordIdentify.Id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		account.DiscordId = discordIdentify.Id
		err = account.Create(c)
	}
	if err != nil {
		redirectErrPage(c)
		slog.Error("Can not ensure an account", "err", err)
		return
	}

	// Issue a jwt for the user
	jwtToken, err := generateJwt(strconv.FormatUint(uint64(account.Id), 10))
	if err != nil {
		redirectErrPage(c)
		slog.Error("Failed to generate jwt", "err", err)
		return
	}

	setCookie(c, jwtCookieName, jwtToken, int(jwtTokenAge/time.Second))
	c.Redirect(302, "/dashboard")
}

func setCookie(c *gin.Context, name string, value string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, value, maxAge, "/", "", config.Conf.IsHttps, true)
}

func handleLogout(c *gin.Context) {
	setCookie(c, jwtCookieName, "", -1)
	c.Redirect(302, "/")
}

func generateJwt(userId string) (string, error) {
	claims := claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtTokenAge)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.Conf.JwtSecret)
}

func generateState() (string, error) {
	const stateLen = 32
	b := make([]byte, stateLen/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func redirectErrPage(c *gin.Context, errs ...string) {
	url := url.URL{Path: "/signin"}
	query := url.Query()
	if len(errs) == 0 {
		errs = []string{"エラーが発生しました。もう一度やり直してください。"}
	}
	var builder strings.Builder
	for _, s := range errs {
		builder.WriteString(s)
	}
	query.Set("err", builder.String())
	url.RawQuery = query.Encode()
	c.Redirect(302, url.String())
}
