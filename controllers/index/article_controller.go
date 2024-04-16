package index

import (
	"fmt"
	"html/template"
	"net/http"
	"project/config"
	"project/helpers"
	"project/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
}

func (c *ArticleController) Search(ct *gin.Context) {
	var err error
	var totalCount, totalPage int
	keyword := ct.Query("keyword")

	pageIndex := ct.Query("pno")
	iPageIndex, _ := strconv.Atoi(pageIndex)

	earticle := models.EArticle{}
	articleList := make([]*models.EArticle, config.PAGE_SIZE)
	articleList, totalCount, totalPage, err = earticle.GetArticleList(iPageIndex, config.PAGE_SIZE, keyword)
	var resultPage helpers.ResultPage
	if err != nil {
		fmt.Println(err)
		resultPage.FrontkendErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	fmt.Println(articleList, totalCount, totalPage)
	resultPage.RenderHtml(ct, http.StatusOK, "article/search", gin.H{
		"articleList": articleList,
		"totalCount":  totalCount,
		"totalPage":   totalPage,
		"keyword":     keyword,
	})
}

func (c *ArticleController) Posts(ct *gin.Context) {
	var err error
	var totalCount, totalPage int
	cateId := ct.Query("cateId")
	iCateId, _ := strconv.Atoi(cateId)

	tagId := ct.Query("tagId")
	iTagId, _ := strconv.Atoi(tagId)

	pageIndex := ct.Query("pno")
	iPageIndex, _ := strconv.Atoi(pageIndex)

	article := models.Article{}
	articleList := make([]*models.Article, config.PAGE_SIZE)
	articleList, totalCount, totalPage, err = article.GetArticleList(iPageIndex, config.PAGE_SIZE, iCateId, iTagId)
	var resultPage helpers.ResultPage
	if err != nil {
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.RenderHtml(ct, http.StatusOK, "article/posts", gin.H{
		"articleList": articleList,
		"totalCount":  totalCount,
		"totalPage":   totalPage,
	})
}

func (c *ArticleController) PostsDetail(ct *gin.Context) {
	var err error
	seoLink := ct.Param("seoLink")
	article := models.Article{}
	err = article.GetOneArticleBySeoLink(seoLink, &article)
	var resultPage helpers.ResultPage
	if err != nil {
		resultPage.FrontkendErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	var oneCategories models.Categories
	err = oneCategories.GetOneCategoryById(article.Categories, &oneCategories)
	if err != nil {
		resultPage.FrontkendErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	resultPage.RenderHtml(ct, http.StatusOK, "article/detail", gin.H{
		"article":       article,
		"content":       template.HTML(article.Content),
		"oneCategories": oneCategories,
	})
}

func (c *ArticleController) About(ct *gin.Context) {
	var err error
	seoLink := "about"
	article := models.Article{}
	err = article.GetOneArticleBySeoLink(seoLink, &article)
	var resultPage helpers.ResultPage
	if err != nil {
		resultPage.FrontkendErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	resultPage.RenderHtml(ct, http.StatusOK, "article/detail", gin.H{
		"article": article,
		"content": template.HTML(article.Content),
	})
}

func (c *ArticleController) Categories(ct *gin.Context) {
	var err error
	categorise := models.Categories{}
	var categoriesList []*models.Categories
	categoriesList, err = categorise.GetAllCategories()
	resultPage := helpers.ResultPage{}
	if err != nil {
		resultPage.FrontkendErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	article := models.Article{}
	newCategoriesList := make(map[int]map[string]interface{})
	var articleList []*models.Article
	for key, perCategory := range categoriesList {
		articleList, _, _, err = article.GetArticleList(1, 5, perCategory.ID, 0)

		newPerCategory := make(map[string]interface{})
		newPerCategory["Name"] = perCategory.Name
		newPerCategory["ID"] = perCategory.ID
		newPerCategory["ArticleList"] = articleList

		newCategoriesList[key] = newPerCategory

	}

	fmt.Printf("%+v", newCategoriesList)

	resultPage.RenderHtml(ct, http.StatusOK, "article/categories", gin.H{
		"categoriesList": newCategoriesList,
	})
}

func (c *ArticleController) Tags(ct *gin.Context) {
	var err error
	tags := models.Tags{}
	var tagsList []*models.Tags
	tagsList, err = tags.GetAllTags()
	var resultPage helpers.ResultPage
	if err != nil {
		resultPage.FrontkendErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	article := models.Article{}
	var iTotal int
	newTagsList := make(map[int]map[string]interface{})
	for key, perTag := range tagsList {
		fmt.Printf("%+v %+v %+v", key, perTag, perTag.ID)
		iTotal, err = article.GetArticleCountBytag(perTag.ID)

		newPerTag := make(map[string]interface{})

		newPerTag["Name"] = perTag.Name
		newPerTag["ID"] = perTag.ID
		newPerTag["Total"] = iTotal
		newTagsList[key] = newPerTag
	}

	fmt.Println(newTagsList)

	resultPage.RenderHtml(ct, http.StatusOK, "article/tags", gin.H{
		"tagsList": newTagsList,
	})
}
