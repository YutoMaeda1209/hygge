package api

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/YutoMaeda1209/hygge/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

const stateCookieName = "oauth_state"
const stateCookieAge = time.Minute * 5 // 5 minutes
const jwtCookieName = "jwt"
const jwtTokenAge = time.Hour * 24 * 7 // 7 days

type authSettings struct {
	oauth        *oauth2.Config
	secureCookie bool
}

type claims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

var settings authSettings
var jwtSecret []byte

func auth() {
	// Get env variables
	clientId := os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH2_CLIENT_SECRET")
	redirectUrl := os.Getenv("OAUTH2_REDIRECT_URL")
	secureCookie, err := strconv.ParseBool(os.Getenv("IS_HTTPS"))
	if err != nil {
		slog.Error("Failed to convert the environment variable IS_HTTPS. Set it this to either true or false.")
		os.Exit(1)
	}
	jwtSecretEnv := os.Getenv("JWT_SECRET")

	if clientId == "" || clientSecret == "" || redirectUrl == "" || jwtSecretEnv == "" {
		slog.Error("Authentication environment variables are not set.")
		os.Exit(1)
	}

	settings = authSettings{
		oauth: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  redirectUrl,
			Scopes:       []string{"identify"},
			Endpoint:     endpoints.Discord,
		},
		secureCookie: secureCookie,
	}
	jwtSecret = []byte(jwtSecretEnv)
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

	url := settings.oauth.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(302, url)
}

func handleCallback(c *gin.Context) {
	// Verify and discard state data
	state, err := c.Cookie(stateCookieName)
	if err != nil || c.Query("state") != state {
		c.String(http.StatusBadRequest, "Invalid state")
		slog.Error("Invalid state", "err", err)
		return
	}
	setCookie(c, stateCookieName, "", -1)

	// Verify oauth2 authorization code
	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Missing code")
		slog.Error("Missing code")
		return
	}

	// Get token from discord
	token, err := settings.oauth.Exchange(c.Request.Context(), code)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to exchange token")
		slog.Error("Failed to exchange token", "err", err)
		return
	}

	// Fetch the discord user tied to the token
	client := settings.oauth.Client(c.Request.Context(), token)
	user, err := model.FetchDiscordUser(client)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch discord user")
		slog.Error("Failed to fetch discord user", "err", err)
		return
	}

	// Issue a jwt for the user
	jwtToken, err := generateJwt(user.Id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate jwt")
		slog.Error("Failed to generate jwt", "err", err)
		return
	}

	setCookie(c, jwtCookieName, jwtToken, int(jwtTokenAge/time.Second))
	c.Status(http.StatusOK)
}

func setCookie(ctx *gin.Context, name string, value string, maxAge int) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(name, value, maxAge, "/", "", settings.secureCookie, true)
}

func generateJwt(userId string) (string, error) {
	claims := claims{
		Name: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtTokenAge)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
