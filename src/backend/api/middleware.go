package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(jwtCookieName)
		if err != nil || tokenString == "" {
			c.String(http.StatusUnauthorized, "Missing authorization token")
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.String(http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*claims); ok {
			if claims.Name != "" {
				c.String(http.StatusUnauthorized, "Invalid username")
				c.Abort()
				return
			}
			c.Next()
		}
	}
}
