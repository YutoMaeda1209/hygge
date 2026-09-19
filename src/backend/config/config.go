package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	OAuth2ClientId     string
	OAuth2ClientSecret string
	OAuth2RedirectUrl  string
	IsHttps            bool
	JwtSecret          string
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

	Conf.OAuth2ClientId = clientId
	Conf.OAuth2ClientSecret = clientSecret
	Conf.OAuth2RedirectUrl = redirectUrl
	Conf.IsHttps = isHttps
	Conf.JwtSecret = jwtSecret
	Conf.DatabaseUrl = dbUrl

	return nil
}
