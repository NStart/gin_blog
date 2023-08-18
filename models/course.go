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

func (courseChapter *CourseChapter) GetChpaterList(courseId, pageIndex, pageSize int) ([]*CourseChapter, int, int, error) {
	var err error
	var chapterList []*CourseChapter
	db := GetDB()

	var totalCount int64
	var iTotalCount int
	err = db.Model(CourseChapter{}).Where("course_id = ?", courseId).Count(&totalCount).Error
	iTotalCount = int(totalCount)

	totalPage := int(math.Ceil(float64(iTotalCount) / float64(pageSize)))

	if err != nil {
		return nil, iTotalCount, totalPage, err
	}

	fmt.Println(courseId)
	err = db.Where("course_id = ?", courseId).Order("created_at desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&chapterList).Error
	fmt.Println(err, chapterList)
	if err != nil {
		return nil, iTotalCount, totalPage, err
	}
	fmt.Printf("%+v %+v", chapterList, pageSize)
	return chapterList, iTotalCount, totalPage, nil
}

func (courseArticle *CourseArticle) GetArticleList(courseId, chapterId, pageIndex, pageSize int) ([]*CourseArticle, int, int, error) {
	var err error
	var articleList []*CourseArticle
	db := GetDB()

	var totalCount int64
	var iTotalCount int
	err = db.Model(CourseArticle{}).Where("course_id = ? and chapter_id = ? ", courseId, chapterId).Count(&totalCount).Error
	iTotalCount = int(totalCount)

	totalPage := int(math.Ceil(float64(iTotalCount) / float64(pageSize)))

	if err != nil {
		return nil, iTotalCount, totalPage, err
	}

	fmt.Println(courseId)
	err = db.Where("course_id = ? and chapter_id = ? ", courseId, chapterId).Order("created_at desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&articleList).Error
	fmt.Println(err, articleList)
	if err != nil {
		return nil, iTotalCount, totalPage, err
	}
	fmt.Printf("%+v %+v", articleList, pageSize)
	return articleList, iTotalCount, totalPage, nil
}

func (Course *Course) GetOneCourseById(id int, oneCourse *Course) error {
	var err error
	db := GetDB()
	err = db.Where("id = ?", id).Take(&oneCourse).Error
	fmt.Printf("%+v", oneCourse)
	return err
}

func (courseArticle *CourseArticle) GetOneArticleById(id int, oneArticle *CourseArticle) error {
	var err error
	db := GetDB()
	err = db.Where("id = ?", id).Take(&courseArticle).Error
	fmt.Printf("%+v", courseArticle)
	return err
}
