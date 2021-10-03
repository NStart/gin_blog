package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/config"
	"strings"
)

func CheckLogin() gin.HandlerFunc {
	return func(ct *gin.Context) {
		session := sessions.Default(ct)
		fullPath := ct.Request.URL.Path
		fmt.Println(fullPath)

		inWhiteList := false
		whiteList := []string{
			"ikebackend/login",
			"ikebackend/dologin",
			"ikebackend/logout",
		}
		for _, perPath := range whiteList {
			if strings.Contains(fullPath, perPath) {
				inWhiteList = true
			}
		}

		if session.Get(config.LOGIN_SESSION_KEY) == nil && strings.Contains(fullPath, "ikebackend") && !inWhiteList {
			ct.Redirect(http.StatusFound, "/ikebackend/login")
			ct.Next()
		}
		ct.Next()
	}
}
