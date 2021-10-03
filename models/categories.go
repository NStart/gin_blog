package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"project/helpers"
)

type Categories struct {
	BaseModel
	Name string `gorm:"type:varchar(100);unique;not null"`
}

func (categories *Categories) GetCategoriesList(pageIndex, pageSize int) ([]*Categories, int, int, error) {
	var err error
	var categoriesList []*Categories
	db := GetDB()

	var totalCount int64
	var iTotalCount int
	err = db.Model(Categories{}).Count(&totalCount).Error
	iTotalCount = int(totalCount)

	totalPage := int(math.Ceil(float64(iTotalCount) / float64(pageSize)))

	if err != nil{
		return nil, iTotalCount, totalPage, err
	}

	err = db.Offset((pageIndex-1)*pageSize).Limit(pageSize).Find(&categoriesList).Error
	if err != nil{
		return nil, iTotalCount, totalPage, err
	}
	fmt.Printf("%+v %+v", categoriesList, pageSize)
	return categoriesList, iTotalCount, totalPage, nil
}

func (categories *Categories) GetAllCategories() (categoriesList []*Categories, err error) {
	db := GetDB()
	err = db.Find(&categoriesList).Error
	if err != nil{
		return nil, err
	}
	return categoriesList, nil
}

func (categories *Categories) CheckCategoryExist(name string) (bool, error) {
	var err error
	db := GetDB()
	//categories := models.Categories{}
	err = db.Where("name = ?", name).Take(&categories).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("查询不到数据")
		return false,nil
	} else if err != nil {
		//如果err不等于record not found错误，又不等于nil，那说明sql执行失败了。
		fmt.Println("查询失败", err)
		return false,err
	}
	fmt.Println(categories)
	return true,nil
}

func (categories *Categories) GetOneCategoryById(id int, oneCategory *Categories) error {
	var err error
	db := GetDB()
	err = db.Where("id = ?", id).Take(&oneCategory).Error
	fmt.Printf("%+v", oneCategory)
	return err
}

func (categories *Categories) AddCategory(categoryData *Categories) bool {
	var err error
	db := GetDB()
	if err = db.Create(&categoryData).Error; err != nil {
		fmt.Println("插入失败", err)
		return false
	}
	return true
}

func (categories *Categories) EditCategory(id int, categoryData *Categories) bool {
	var err error
	db := GetDB()
	data := make(map[string]interface{})
	//data["update_at"] =
	data["name"] = categoryData.Name

	err = db.Model(&Categories{}).Where("id = ?", id).Updates(data).Error
	if err != nil{
		return false
	}
	return true
}

func (categories *Categories) CheckCategoryIsUsed(cateId int) (bool,error) {
	var err error
	db := GetDB()
	article := Article{}
	err = db.Where("categories = ?", cateId).Take(&article).Error

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

func (categories *Categories) DeleteCategoryById(id int) (bool, error) {
	db := GetDB()

	var err error
	var isUsed bool
	isUsed, err = categories.CheckCategoryIsUsed(id)
	fmt.Printf("is used: %+v is error %+v:", isUsed, err)

	if err != nil{
		return false,errors.New(helpers.LAN_DB_ERROR)
	}
	if isUsed{
		return false,errors.New(helpers.LAN_CATEGORIES_IS_USED)
	}

	err = db.Where("id = ?", id).Delete(&Categories{}).Error
	if err != nil{
		return false,errors.New(helpers.LAN_DB_ERROR)
	}
	return true, nil
}

