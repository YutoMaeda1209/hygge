package conf

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Conf struct {
	OAuth2ClientId     string
	OAuth2ClientSecret string
	OAuth2RedirectUrl  string
	IsHttps            bool
	JwtSecret          string
	DatabaseUrl        string
}

var conf Conf

func LoadConf() error {
	err := godotenv.Load()
	if err != nil {
		slog.Info("dotenv file is not loaded.")
	}

	clientId := os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH2_CLIENT_SECRET")
	redirectUrl := os.Getenv("OAUTH2_REDIRECT_URL")
	isHttps, err := strconv.ParseBool(os.Getenv("IS_HTTPS"))
	jwtSecret := os.Getenv("JWT_SECRET")
	dbUrl := os.Getenv("DB_URL")

	if clientId == "" || clientSecret == "" || redirectUrl == "" {
		return errors.New("OAuth2 environment variables are not set.")
	} else if err != nil || jwtSecret == "" {
		return errors.New("Jwt environment variables are not set.")
	} else if dbUrl == "" {
		return errors.New("Database environment variables are not set.")
	}

	conf.OAuth2ClientId = clientId
	conf.OAuth2ClientSecret = clientSecret
	conf.OAuth2RedirectUrl = redirectUrl
	conf.IsHttps = isHttps
	conf.JwtSecret = jwtSecret
	conf.DatabaseUrl = dbUrl

	return nil
}
