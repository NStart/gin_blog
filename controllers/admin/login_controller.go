package admin

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/config"
	"project/helpers"
	"project/models"
)

type LoginController struct {

}

func (c *LoginController) Index(ct *gin.Context)  {
	var resultPage helpers.ResultPage
	resultPage.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/login/index", gin.H{

	})
}

func (c *LoginController) DoLogin(ct *gin.Context)  {
	var err error
	name := ct.PostForm("name")
	password := ct.PostForm("password")

	user := models.AdminUser{}
	err = user.GetOneUserByName(name, &user)
	var resultPage helpers.ResultPage
	if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_UER_NOT_EXIST)
		return
	}

	md5Password := md5.Sum([]byte(password+config.SESSION_SECRET))
	sMd5Password := fmt.Sprintf("%x", md5Password)
	if sMd5Password != user.Password {
		resultPage.AdminErrorPage(ct, helpers.LAN_PASSWORD_ERROR)
		return
	}

	session := sessions.Default(ct)
	session.Set(config.LOGIN_SESSION_KEY, user.Name)
	session.Save()

	ct.Redirect(http.StatusMovedPermanently, "/ikebackend/article/index")
}

func (c *LoginController) Logout(ct *gin.Context)  {
	session := sessions.Default(ct)
	session.Delete(config.LOGIN_SESSION_KEY)
	session.Save()

	ct.Redirect(http.StatusFound, "/ikebackend/login")
}