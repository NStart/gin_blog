package index

import (
	"html/template"
	"net/http"
	"project/config"
	"project/helpers"
	"project/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
}

func (c *CourseController) List(ct *gin.Context) {
	var err error
	var totalCount, totalPage int

	pageIndex := ct.Query("pno")
	iPageIndex, _ := strconv.Atoi(pageIndex)

	course := models.Course{}
	courseList := make([]*models.Course, config.PAGE_SIZE)
	courseList, totalCount, totalPage, err = course.GetCourseList(iPageIndex, config.PAGE_SIZE)
	var resultPage helpers.ResultPage
	if err != nil {
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.RenderHtml(ct, http.StatusOK, "course/index", gin.H{
		"courseList": courseList,
		"totalCount": totalCount,
		"totalPage":  totalPage,
	})
}

func (c *CourseController) Posts(ct *gin.Context) {
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

func (c *CourseController) PostsDetail(ct *gin.Context) {
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
