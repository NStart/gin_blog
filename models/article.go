package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
)

type Article struct {
	BaseModel
	SeoLink string `gorm:"type:varchar(100);unique;not null"`
	Categories int `gorm:"not null"`
	Tags int `gorm:"not null"`
	Title string `gorm:"type:varchar(100);not null"`
	Content string `gorm:"not null"`
}

func (article *Article) GetArticleList(pageIndex, pageSize, cateId, tagId int) ([]*Article, int, int, error) {
	var err error
	var articleList []*Article
	db := GetDB()

	var totalCount int64
	var iTotalCount int
	err = db.Model(Article{}).Count(&totalCount).Error
	iTotalCount = int(totalCount)

	totalPage := int(math.Ceil(float64(iTotalCount) / float64(pageSize)))

	if err != nil{
		return nil, iTotalCount, totalPage, err
	}

	if cateId != 0{
		db = db.Where("categories = ?", cateId)
	}
	if tagId != 0{
		db = db.Where("tags = ?", tagId)
	}

	err = db.Order("created_at desc").Offset((pageIndex-1)*pageSize).Limit(pageSize).Find(&articleList).Error
	if err != nil{
		return nil, iTotalCount, totalPage, err
	}
	fmt.Printf("%+v %+v", articleList, pageSize)
	return articleList, iTotalCount, totalPage, nil
}

func (article *Article) CheckArticleExist(seoLink string) (bool, error) {
	var err error
	db := GetDB()
	//article := models.Article{}
	err = db.Where("seo_link = ?", seoLink).Take(&article).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("查询不到数据")
		return false,nil
	} else if err != nil {
		//如果err不等于record not found错误，又不等于nil，那说明sql执行失败了。
		fmt.Println("查询失败", err)
		return false,err
	}
	fmt.Println(article)
	return true,nil
}

func (article *Article) GetArticleCountBytag(tagId int) (int, error) {
	var err error
	db := GetDB()
	var total int64
	db.Model(&article).Where("tags = ?", tagId).Count(&total)
	var iTotal int
	iTotal = int(total)
	fmt.Printf("%+v", iTotal)
	return iTotal, err
}

func (article *Article) GetOneArticleById(id int, oneArticle *Article) error {
	var err error
	db := GetDB()
	err = db.Where("id = ?", id).Take(&oneArticle).Error
	fmt.Printf("%+v", oneArticle)
	return err
}

func (article *Article) GetOneArticleBySeoLink(seoLink string, oneArticle *Article) error {
	var err error
	db := GetDB()
	err = db.Where("seo_link = ?", seoLink).Take(&oneArticle).Error
	fmt.Printf("%+v", oneArticle)
	return err
}

func (article *Article) AddArticle(articleData *Article) bool {
	var err error
	db := GetDB()
	if err = db.Create(&articleData).Error; err != nil {
		fmt.Println("插入失败", err)
		return false
	}
	return true
}

func (article *Article) EditArticle(id int, articleData *Article) bool {
	var err error
	db := GetDB()
	data := make(map[string]interface{})
	//data["update_at"] =
	data["seo_link"] = articleData.SeoLink
	data["categories"] = articleData.Categories
	data["tags"] = articleData.Tags
	data["title"] = articleData.Title
	data["content"] = articleData.Content

	err = db.Model(&Article{}).Where("id = ?", id).Updates(data).Error
	if err != nil{
		return false
	}
	return true
}

func (article *Article) DeleteArticleById(id int) bool {
	db := GetDB()
	err := db.Where("id = ?", id).Delete(&Article{}).Error
	if err != nil{
		return false
	}
	return true
}


