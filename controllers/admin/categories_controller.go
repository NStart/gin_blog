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

type CategoriesController struct {

}

func (c *CategoriesController) Index(ct *gin.Context)  {
	var err error
	var totalCount, totalPage int
	pageIndex := ct.Query("pno")
	iPageIndex, _ := strconv.Atoi(pageIndex)
	categories := models.Categories{}
	categoriesList := make([]*models.Categories, config.PAGE_SIZE)
	categoriesList, totalCount, totalPage, err = categories.GetCategoriesList(iPageIndex, config.PAGE_SIZE)
	var resultPage helpers.ResultPage
	if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/categories/index", gin.H{
		"categoriesList": categoriesList,
		"totalCount": totalCount,
		"totalPage" : totalPage,
	})
}

func (c *CategoriesController) Add(ct *gin.Context)  {
	var resultPage helpers.ResultPage
	resultPage.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/categories/add", gin.H{
		"title": "Main website ",
	})
}

func (c *CategoriesController) Edit(ct *gin.Context)  {
	var err error
	id := ct.Query("id")
	iId, err := strconv.Atoi(id)
	categories := models.Categories{}
	err = categories.GetOneCategoryById(iId, &categories)
	var resultPage helpers.ResultPage
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resultPage.AdminErrorPage(ct, helpers.LAN_NOT_EXIST_RECORD)
		return
	} else if err != nil {
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	//fmt.Println( tags, 123345)
	resultPage.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/categories/edit", gin.H{
		"categories": categories,
	})
}

func (c *CategoriesController) Doadd(ct *gin.Context)  {
	var err error
	name := ct.PostForm("name")

	categories := models.Categories{}
	var isExist bool
	isExist, err = categories.CheckCategoryExist(name)
	var resultPage helpers.ResultPage
	if isExist {
		resultPage.AdminErrorPage(ct, helpers.LAN_EXIST_RECORD)
		return
	}else if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	categories = models.Categories{
		Name: name,
	}

	if !categories.AddCategory(&categories) {
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.AdminSuccPage(ct, helpers.LAN_ADD_SUCC)
	return
}

func (c *CategoriesController) Doedit(ct *gin.Context)  {
	//var err error
	var iId int
	id := ct.PostForm("id")
	iId, _ = strconv.Atoi(id)
	name := ct.PostForm("name")


	categories := models.Categories{
		Name: name,
	}

	var resultPage helpers.ResultPage
	if !categories.EditCategory(iId, &categories){
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.AdminSuccPage(ct, helpers.LAN_EDIT_SUCC)
	return
}

func (c *CategoriesController) Delete(ct *gin.Context)  {
	var err error
	var iId int
	id := ct.Query("id")
	iId, _ = strconv.Atoi(id)

	//fmt.Printf("%+v", iId)
	categories := models.Categories{}
	var JSONResponse helpers.JSONResponse
	var deleteResult bool
	deleteResult, err = categories.DeleteCategoryById(iId)
	if err != nil{
		JSONResponse.FailJSON(ct, err.Error(), nil)
		return
	}
	if !deleteResult{
		JSONResponse.FailJSON(ct, helpers.LAN_DB_ERROR, nil)
		return
	}
	JSONResponse.SuccJSON(ct, helpers.LAN_DELETE_SUCC, nil)
	return
}