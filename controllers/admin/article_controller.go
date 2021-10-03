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

type ArticleController struct {

}

func (c *ArticleController) Index(ct *gin.Context)  {
	var err error
	var totalCount, totalPage int
	pageIndex := ct.Query("pno")
	iPageIndex, _ := strconv.Atoi(pageIndex)
	article := models.Article{}
	articleList := make([]*models.Article, config.PAGE_SIZE)
	articleList, totalCount, totalPage, err = article.GetArticleList(iPageIndex, config.PAGE_SIZE, 0, 0)
	var resultPage helpers.ResultPage
	if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/article/index", gin.H{
		"articleList": articleList,
		"totalCount": totalCount,
		"totalPage" : totalPage,
	})
}

func (c *ArticleController) Add(ct *gin.Context)  {
	var err error
	tags := models.Tags{}
	var tagsList []*models.Tags
	tagsList, err = tags.GetAllTags()
	var resultPage helpers.ResultPage
	if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	categories := models.Categories{}
	var categoriesList []*models.Categories
	categoriesList, err = categories.GetAllCategories()
	if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	resultPage.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/article/add", gin.H{
		"tagsList": tagsList,
		"categoriesList": categoriesList,
	})
}

func (c *ArticleController) Edit(ct *gin.Context)  {
	var err error
	id := ct.Query("id")
	iId, err := strconv.Atoi(id)
	article := models.Article{}
	err = article.GetOneArticleById(iId, &article)
	var resultPage helpers.ResultPage
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resultPage.AdminErrorPage(ct, helpers.LAN_NOT_EXIST_RECORD)
		return
	} else if err != nil {
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	tags := models.Tags{}
	var tagsList []*models.Tags
	tagsList, err = tags.GetAllTags()
	if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	categories := models.Categories{}
	var categoriesList []*models.Categories
	categoriesList, err = categories.GetAllCategories()
	if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	//fmt.Println( article, 123345)
	resultPage.RenderHtml(ct, http.StatusOK, config.ADMIN_DIR+"/article/edit", gin.H{
		"article": article,
		"tagsList": tagsList,
		"categoriesList": categoriesList,
	})
}

func (c *ArticleController) Doadd(ct *gin.Context)  {
	var err error
	var iCategories, iTags int
	seoLink := ct.PostForm("seo_link")
	categories := ct.PostForm("categories")
	iCategories,err = strconv.Atoi(categories)
	tags := ct.PostForm("tags")
	iTags,err = strconv.Atoi(tags)
	title := ct.PostForm("title")
	content := ct.PostForm("content")

	article := models.Article{}
	var isExist bool
	isExist, err = article.CheckArticleExist(seoLink)
	var resultPage helpers.ResultPage
	if isExist {
		resultPage.AdminErrorPage(ct, helpers.LAN_EXIST_RECORD)
		return
	}else if err != nil{
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	article = models.Article{
		Title:      title,
		SeoLink:    seoLink,
		Categories: iCategories,
		Tags:       iTags,
		Content:    content,
	}

	if !article.AddArticle(&article) {
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.AdminSuccPage(ct, helpers.LAN_ADD_SUCC)
	return
}

func (c *ArticleController) Doedit(ct *gin.Context)  {
	//var err error
	var iId, iCategories, iTags int
	id := ct.PostForm("id")
	iId, _ = strconv.Atoi(id)
	seoLink := ct.PostForm("seo_link")
	categories := ct.PostForm("categories")
	iCategories,_ = strconv.Atoi(categories)
	tags := ct.PostForm("tags")
	iTags, _ = strconv.Atoi(tags)
	title := ct.PostForm("title")
	content := ct.PostForm("content")

	article := models.Article{
		SeoLink: seoLink,
		Categories: iCategories,
		Tags: iTags,
		Title: title,
		Content: content,
	}

	var resultPage helpers.ResultPage
	if !article.EditArticle(iId, &article){
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.AdminSuccPage(ct, helpers.LAN_EDIT_SUCC)
	return
}

func (c *ArticleController) Delete(ct *gin.Context)  {
	var iId int
	id := ct.Query("id")
	iId, _ = strconv.Atoi(id)

	//fmt.Printf("%+v", iId)
	article := models.Article{}
	var JSONResponse helpers.JSONResponse
	if !article.DeleteArticleById(iId){
		JSONResponse.FailJSON(ct, helpers.LAN_DB_ERROR, nil)
		return
	}
	JSONResponse.SuccJSON(ct, helpers.LAN_DELETE_SUCC, nil)
	return
}