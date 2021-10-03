package index

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/helpers"
)

type IndexController struct {

}



func (c *IndexController) Index(ct *gin.Context) {
	var resultPage helpers.ResultPage
	resultPage.RenderHtml(ct, http.StatusOK, "index/index", gin.H{
		"title": "Main website",
	})
}

func (c *IndexController) About(ct *gin.Context) {
	var resultPage helpers.ResultPage
	resultPage.RenderHtml(ct, http.StatusOK, "index/about", gin.H{
		"title": "Main website",
	})
}
