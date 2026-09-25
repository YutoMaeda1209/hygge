package api

import (
	"errors"
	"net/http"

	"github.com/YutoMaeda1209/hygge/config"
	"github.com/YutoMaeda1209/hygge/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(jwtCookieName)
		if err != nil || tokenString == "" {
			c.String(http.StatusUnauthorized, "Missing authorization token")
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (any, error) {
			return config.Conf.JwtSecret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || !token.Valid {
			c.String(http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*claims)
		if !ok {
			c.String(http.StatusInternalServerError, "")
			c.Abort()
			return
		}

		account := model.Account{}
		err = account.Find(c.Request.Context(), "id = ?", claims.UserId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusUnauthorized, "Invalid username")
			c.Abort()
			return
		} else if err != nil {
			c.String(http.StatusInternalServerError, "")
			c.Abort()
			return
		}

		c.Set(accountContextKey, account)
		c.Next()
	}
}
