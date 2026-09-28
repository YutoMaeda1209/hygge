package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type Config struct {
	OAuth2Config *oauth2.Config
	BotToken     string
	IsHttps      bool
	DatabaseUrl  string
	TrustProxyIp []string
}

var Conf Config

// LoadConf loads settings from the environment (and .env if present), validates them,
// and stores them in Conf. It returns an error if a required variable is missing or invalid.
func LoadConf() error {
	err := godotenv.Load()
	if err != nil {
		slog.Info("dotenv file is not loaded.")
	}

	clientId := os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH2_CLIENT_SECRET")
	redirectUrl := os.Getenv("OAUTH2_REDIRECT_URL")
	isHttps, parseBoolErr := strconv.ParseBool(os.Getenv("IS_HTTPS"))
	dbUrl := os.Getenv("DB_URL")
	trustProxyIp := os.Getenv("TRUST_PROXY_IP")
	botToken := os.Getenv("DISCORD_BOT_TOKEN")

	if clientId == "" || clientSecret == "" || redirectUrl == "" {
		return errors.New("OAuth2 environment variables are not set")
	} else if parseBoolErr != nil {
		return errors.New("IS_HTTPS must be set to a boolean value")
	} else if botToken == "" {
		return errors.New("DISCORD_BOT_TOKEN is not set")
	} else if dbUrl == "" {
		return errors.New("the database URL has not been set; visit https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING-KEYWORD-VALUE to configure it")
	}

	Conf.OAuth2Config = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"identify"},
		Endpoint:     endpoints.Discord,
	}
	Conf.BotToken = botToken
	Conf.IsHttps = isHttps
	Conf.DatabaseUrl = dbUrl

	// Format validation is left to gin's SetTrustedProxies.
	var ips []string
	for ip := range strings.SplitSeq(trustProxyIp, ",") {
		if ip = strings.TrimSpace(ip); ip != "" {
			ips = append(ips, ip)
		}
	}
	Conf.TrustProxyIp = ips

	return nil
}
