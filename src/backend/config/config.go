package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type Config struct {
	OAuth2ClientId     string
	OAuth2ClientSecret string
	OAuth2RedirectUrl  string
	OAuth2Config       *oauth2.Config
	IsHttps            bool
	JwtSecret          []byte
	DatabaseUrl        string
}

var Conf Config

func LoadConf() error {
	err := godotenv.Load()
	if err != nil {
		slog.Info("dotenv file is not loaded.")
	}

	clientId := os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH2_CLIENT_SECRET")
	redirectUrl := os.Getenv("OAUTH2_REDIRECT_URL")
	isHttps, parseBoolErr := strconv.ParseBool(os.Getenv("IS_HTTPS"))
	jwtSecret := os.Getenv("JWT_SECRET")
	dbUrl := os.Getenv("DB_URL")

	if clientId == "" || clientSecret == "" || redirectUrl == "" {
		return errors.New("OAuth2 environment variables are not set.")
	} else if parseBoolErr != nil || jwtSecret == "" {
		return errors.New("Jwt environment variables are not set.")
	} else if dbUrl == "" {
		return errors.New("The database URL has not been set. Visit https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING-KEYWORD-VALUE to configure it.")
	}

	Conf.OAuth2ClientId = clientId
	Conf.OAuth2ClientSecret = clientSecret
	Conf.OAuth2RedirectUrl = redirectUrl
	Conf.OAuth2Config = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"identify"},
		Endpoint:     endpoints.Discord,
	}
	Conf.IsHttps = isHttps
	Conf.JwtSecret = []byte(jwtSecret)
	Conf.DatabaseUrl = dbUrl

	return nil
}
