package models

import (
	"fmt"
	"math"
)

type Course struct {
	BaseModel
	CourseName  string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"not null"`
}

type CourseChapter struct {
	BaseModel
	CourseId    int64  `gorm:"not null"`
	ChapterName string `gorm:"not null"`
}

type CourseArticle struct {
	BaseModel
	CourseId       int64  `gorm:"not null"`
	ArticleName    string `gorm:"not null"`
	ArticleContent string `gorm:"not null"`
}

func (course *Course) GetCourseList(pageIndex, pageSize int) ([]*Course, int, int, error) {
	var err error
	var courseList []*Course
	db := GetDB()

	var totalCount int64
	var iTotalCount int
	err = db.Model(Course{}).Count(&totalCount).Error
	iTotalCount = int(totalCount)

	totalPage := int(math.Ceil(float64(iTotalCount) / float64(pageSize)))

	if err != nil {
		return nil, iTotalCount, totalPage, err
	}

	err = db.Order("created_at desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&courseList).Error
	if err != nil {
		return nil, iTotalCount, totalPage, err
	}
	fmt.Printf("%+v %+v", courseList, pageSize)
	return courseList, iTotalCount, totalPage, nil
}

func (Course *Course) GetOneCourseById(id int, oneCourse *Course) error {
	var err error
	db := GetDB()
	err = db.Where("id = ?", id).Take(&oneCourse).Error
	fmt.Printf("%+v", oneCourse)
	return err
}

func (c *CourseChapter) GetChapterList(pageIndex, pageSize int) ([]*CourseChapter, int, int, error) {
	var err error
	var chapterList []*CourseChapter
	db := GetDB()

	var totalCount int64
	var iTotalCount int
	err = db.Model(Course{}).Count(&totalCount).Error
	iTotalCount = int(totalCount)

	totalPage := int(math.Ceil(float64(iTotalCount) / float64(pageSize)))

	if err != nil {
		return nil, iTotalCount, totalPage, err
	}

	err = db.Order("created_at desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&chapterList).Error
	if err != nil {
		return nil, iTotalCount, totalPage, err
	}
	fmt.Printf("%+v %+v", chapterList, pageSize)
	return chapterList, iTotalCount, totalPage, nil
}

func (c *CourseChapter) GetOneChapterById(id int, oneChapter *Course) error {
	var err error
	db := GetDB()
	err = db.Where("id = ?", id).Take(&oneChapter).Error
	fmt.Printf("%+v", oneChapter)
	return err
}

func (a *CourseArticle) GetArticleList(pageIndex, pageSize, courseId int) ([]*CourseArticle, int, int, error) {
	var err error
	var articleList []*CourseArticle
	db := GetDB()

	var totalCount int64
	var iTotalCount int
	err = db.Model(Course{}).Count(&totalCount).Error
	iTotalCount = int(totalCount)

	totalPage := int(math.Ceil(float64(iTotalCount) / float64(pageSize)))

	if err != nil {
		return nil, iTotalCount, totalPage, err
	}

	err = db.Order("created_at desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&articleList).Error
	if err != nil {
		return nil, iTotalCount, totalPage, err
	}
	fmt.Printf("%+v %+v", articleList, pageSize)
	return articleList, iTotalCount, totalPage, nil
}

func (a *CourseArticle) GetOneArticleById(id int, oneArticle *Course) error {
	var err error
	db := GetDB()
	err = db.Where("id = ?", id).Take(&oneArticle).Error
	fmt.Printf("%+v", oneArticle)
	return err
}
