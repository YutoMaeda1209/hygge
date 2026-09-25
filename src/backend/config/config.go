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

// HS256 requires a key at least as long as the hash output (RFC 7518 3.2).
const minJwtSecretLen = 32

type Config struct {
	OAuth2Config *oauth2.Config
	IsHttps      bool
	JwtSecret    []byte
	DatabaseUrl  string
	TrustProxyIp []string
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
	trustProxyIp := os.Getenv("TRUST_PROXY_IP")

	if clientId == "" || clientSecret == "" || redirectUrl == "" {
		return errors.New("OAuth2 environment variables are not set")
	} else if parseBoolErr != nil {
		return errors.New("IS_HTTPS must be set to a boolean value")
	} else if len(jwtSecret) < minJwtSecretLen {
		return errors.New("JWT_SECRET must be at least 32 bytes")
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
	Conf.IsHttps = isHttps
	Conf.JwtSecret = []byte(jwtSecret)
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
