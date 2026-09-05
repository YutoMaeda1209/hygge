package api

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

const stateCookieName = "oauth_state"
const stateCookieAge = int(time.Minute * 5 / time.Second) // 5 minutes

type authSettings struct {
	oauth        *oauth2.Config
	secureCookie bool
}

var settings authSettings

func auth() {
	// Get env variables
	clientId := os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH2_CLIENT_SECRET")
	redirectUrl := os.Getenv("OAUTH2_REDIRECT_URL")
	secureCookie, err := strconv.ParseBool(os.Getenv("IS_HTTPS"))
	if err != nil {
		log.Fatalln("Failed to convert the environment variable IS_HTTPS. Set it this to either true or false.")
	}

	if clientId == "" || clientSecret == "" || redirectUrl == "" {
		log.Fatalln("OAuth2 environment variables are not set.")
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
}

func handleLogin(c *gin.Context) {
	// Generate and store state data
	state, err := generateState()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate state")
		return
	}
	c.SetCookie(stateCookieName, state, stateCookieAge, "/", "", settings.secureCookie, true)

	url := settings.oauth.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(302, url)
}

func handleCallback(c *gin.Context) {
	// Verify and discard state data
	state, err := c.Cookie(stateCookieName)
	if err != nil || c.Query("state") != state {
		c.String(http.StatusBadRequest, "Invalid state")
		return
	}
	c.SetCookie(stateCookieName, "", -1, "/", "", settings.secureCookie, true)

	// Verify oauth2 authorization code
	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Missing code")
		return
	}

	// Get token from discord
	token, err := settings.oauth.Exchange(c.Request.Context(), code)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to exchange token: %v", err)
		return
	}

	_ = token
	c.Status(http.StatusOK)
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
