package api

import (
	"net/http"

	"github.com/YutoMaeda1209/hygge/config"
	"github.com/gin-gonic/gin"
)

func setCookie(c *gin.Context, name string, value string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, value, maxAge, "/", "", config.Conf.IsHttps, true)
}

func removeCookie(c *gin.Context, name string) {
	setCookie(c, name, "", -1)
}
