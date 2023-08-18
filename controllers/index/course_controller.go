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

func (c *CourseController) Courses(ct *gin.Context) {
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

func (c *CourseController) ChaperList(ct *gin.Context) {
	var err error
	var totalCount, totalPage int

	courseId := ct.Param("courseId")
	iCourseId, _ := strconv.Atoi(courseId)
	pageIndex := ct.Query("pno")
	iPageIndex, _ := strconv.Atoi(pageIndex)

	courseChapter := models.CourseChapter{}
	chapterList := make([]*models.CourseChapter, config.PAGE_SIZE)
	chapterList, totalCount, totalPage, err = courseChapter.GetChpaterList(iCourseId, iPageIndex, config.PAGE_SIZE)
	var resultPage helpers.ResultPage
	if err != nil {
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.RenderHtml(ct, http.StatusOK, "course/chapter", gin.H{
		"chapterList": chapterList,
		"totalCount":  totalCount,
		"totalPage":   totalPage,
	})
}

func (c *CourseController) CourseArticleList(ct *gin.Context) {
	var err error
	var totalCount, totalPage int

	courseId := ct.Param("courseId")
	iCourseId, _ := strconv.Atoi(courseId)
	chapterId := ct.Param("chapterId")
	iChapterId, _ := strconv.Atoi(chapterId)
	pageIndex := ct.Query("pno")
	iPageIndex, _ := strconv.Atoi(pageIndex)

	courseArticle := models.CourseArticle{}
	aiticleList := make([]*models.CourseArticle, config.PAGE_SIZE)
	aiticleList, totalCount, totalPage, err = courseArticle.GetArticleList(iCourseId, iChapterId, iPageIndex, config.PAGE_SIZE)
	var resultPage helpers.ResultPage
	if err != nil {
		resultPage.AdminErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}
	resultPage.RenderHtml(ct, http.StatusOK, "course/article", gin.H{
		"aiticleList": aiticleList,
		"totalCount":  totalCount,
		"totalPage":   totalPage,
	})
}

func (c *CourseController) ArticleDetail(ct *gin.Context) {
	var err error
	articleId := ct.Param("articleId")
	iArticleId, _ := strconv.Atoi(articleId)
	article := models.CourseArticle{}
	err = article.GetOneArticleById(iArticleId, &article)
	var resultPage helpers.ResultPage
	if err != nil {
		resultPage.FrontkendErrorPage(ct, helpers.LAN_DB_ERROR)
		return
	}

	resultPage.RenderHtml(ct, http.StatusOK, "course/article_detail", gin.H{
		"article": article,
		"content": template.HTML(article.ArticleContent),
	})
}
