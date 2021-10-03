package helpers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/config"
)

type ResultPage struct {

}

func (h *ResultPage) FrontkendErrorPage(ct *gin.Context, errorMsg string)  {
	h.RenderHtml(ct, http.StatusOK, "common/error", gin.H{
		"errorMsg": errorMsg,
	})
}

func (h *ResultPage) FrontkendSuccPage(ct *gin.Context, succMsg string)  {
	h.RenderHtml(ct, http.StatusOK, "common/succ", gin.H{
		"succMsg": succMsg,
	})
}

func (h *ResultPage) AdminErrorPage(ct *gin.Context, errorMsg string)  {
	h.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/common/error", gin.H{
		"errorMsg": errorMsg,
	})
}

func (h *ResultPage) AdminSuccPage(ct *gin.Context, succMsg string)  {
	h.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/common/succ", gin.H{
		"succMsg": succMsg,
	})
}

func (h *ResultPage) RenderHtml(ct *gin.Context, code int, name string, obj gin.H)  {
	session := sessions.Default(ct)
	adminSession := session.Get(config.LOGIN_SESSION_KEY)
	sAdminSession := fmt.Sprintf("%+v", adminSession)
	if adminSession == nil {
		obj["adminIsLogin"] = 0
		obj["adminLoginName"] = ""
	}else{
		obj["adminIsLogin"] = 1
		obj["adminLoginName"] = sAdminSession
	}
	obj["domainName"] = config.DOMAIN_NAME
	fmt.Println(obj)
	ct.HTML(code, name, obj)
}

