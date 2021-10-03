package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"project/config"
	"project/helpers"
	"project/models"
	"strconv"
)

type TagsController struct {

}

func (c *TagsController) Index(ct *gin.Context)  {
	var err error
	var totalCount, totalPage int
	pageIndex := ct.Query("pno")
	iPageIndex, _ := strconv.Atoi(pageIndex)
	tags := models.Tags{}
	tagsList := make([]*models.Tags, config.PAGE_SIZE)
	tagsList, totalCount, totalPage, err = tags.GetTagsList(iPageIndex, config.PAGE_SIZE)
	var resultPage helpers.ResultPage
	if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/tags/index", gin.H{
		"tagsList": tagsList,
		"totalCount": totalCount,
		"totalPage" : totalPage,
	})
}

func (c *TagsController) Add(ct *gin.Context)  {
	var resultPage helpers.ResultPage
	resultPage.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/tags/add", gin.H{
		"title": "Main website ",
	})
}

func (c *TagsController) Edit(ct *gin.Context)  {
	var err error
	id := ct.Query("id")
	iId, err := strconv.Atoi(id)
	tags := models.Tags{}
	err = tags.GetOneTagById(iId, &tags)
	var resultPage helpers.ResultPage
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resultPage.AdminErrorPage(ct, helpers.LAN_NOT_EXIST_RECORD)
		return
	} else if err != nil {
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	//fmt.Println( tags, 123345)
	resultPage.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/tags/edit", gin.H{
		"tags": tags,
	})
}

func (c *TagsController) Doadd(ct *gin.Context)  {
	var err error
	name := ct.PostForm("name")

	tags := models.Tags{}
	var isExist bool
	isExist, err = tags.CheckTagExist(name)
	var resultPage helpers.ResultPage
	if isExist {
		resultPage.AdminErrorPage(ct, helpers.LAN_EXIST_RECORD)
		return
	}else if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	tags = models.Tags{
		Name: name,
	}

	if !tags.AddTag(&tags) {
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.AdminSuccPage(ct, helpers.LAN_ADD_SUCC)
	return
}

func (c *TagsController) Doedit(ct *gin.Context)  {
	//var err error
	var iId int
	id := ct.PostForm("id")
	iId, _ = strconv.Atoi(id)
	name := ct.PostForm("name")


	tags := models.Tags{
		Name: name,
	}

	var resultPage helpers.ResultPage
	if !tags.EditTag(iId, &tags){
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.AdminSuccPage(ct, helpers.LAN_EDIT_SUCC)
	return
}

func (c *TagsController) Delete(ct *gin.Context)  {
	var err error
	var iId int
	id := ct.Query("id")
	iId, _ = strconv.Atoi(id)

	//fmt.Printf("%+v", iId)
	tags := models.Tags{}
	var JSONResponse helpers.JSONResponse
	var deleteResult bool
	deleteResult, err = tags.DeleteTagById(iId)
	if err != nil || !deleteResult{
		JSONResponse.FailJSON(ct, err.Error(), nil)
		return
	}

	JSONResponse.SuccJSON(ct, helpers.LAN_DELETE_SUCC, nil)
	return
}