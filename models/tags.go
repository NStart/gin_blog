package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"project/helpers"
)

type Tags struct {
	BaseModel
	Name string `gorm:"type:varchar(100);unique;not null"`
}

func (tags *Tags) GetTagsList(pageIndex, pageSize int) ([]*Tags, int, int, error) {
	var err error
	var tagsList []*Tags
	db := GetDB()

	var totalCount int64
	var iTotalCount int
	err = db.Model(Tags{}).Count(&totalCount).Error
	iTotalCount = int(totalCount)

	totalPage := int(math.Ceil(float64(iTotalCount) / float64(pageSize)))

	if err != nil{
		return nil, iTotalCount, totalPage, err
	}

	err = db.Offset((pageIndex-1)*pageSize).Limit(pageSize).Find(&tagsList).Error
	if err != nil{
		return nil, iTotalCount, totalPage, err
	}
	fmt.Printf("%+v %+v", tagsList, pageSize)
	return tagsList, iTotalCount, totalPage, nil
}

func (tags *Tags) GetAllTags() (tagsList []*Tags, err error) {
	db := GetDB()
	err = db.Find(&tagsList).Error
	if err != nil{
		return nil, err
	}
	return tagsList, nil
}

func (tags *Tags) CheckTagExist(name string) (bool, error) {
	var err error
	db := GetDB()
	//article := models.Article{}
	err = db.Where("name = ?", name).Take(&tags).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("查询不到数据")
		return false,nil
	} else if err != nil {
		//如果err不等于record not found错误，又不等于nil，那说明sql执行失败了。
		fmt.Println("查询失败", err)
		return false,err
	}
	fmt.Println(tags)
	return true,nil
}

func (tags *Tags) GetOneTagById(id int, oneTag *Tags) error {
	var err error
	db := GetDB()
	err = db.Where("id = ?", id).Take(&oneTag).Error
	fmt.Printf("%+v", oneTag)
	return err
}

func (tags *Tags) AddTag(tagData *Tags) bool {
	var err error
	db := GetDB()
	if err = db.Create(&tagData).Error; err != nil {
		fmt.Println("插入失败", err)
		return false
	}
	return true
}

func (tags *Tags) EditTag(id int, tagData *Tags) bool {
	var err error
	db := GetDB()
	data := make(map[string]interface{})
	//data["update_at"] =
	data["name"] = tagData.Name

	err = db.Model(&Tags{}).Where("id = ?", id).Updates(data).Error
	if err != nil{
		return false
	}
	return true
}

func (tags *Tags) CheckTagIsUsed(tagId int) (bool,error) {
	var err error
	db := GetDB()
	article := Article{}
	err = db.Where("tags = ?", tagId).Take(&article).Error

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

func (tags *Tags) DeleteTagById(id int) (bool, error) {
	db := GetDB()

	var err error
	var isUsed bool
	isUsed, err = tags.CheckTagIsUsed(id)

	if err != nil{
		return false,errors.New(helpers.LAN_DB_ERROR)
	}
	if isUsed{
		return false,errors.New(helpers.LAN_TAG_IS_USED)
	}

	err = db.Where("id = ?", id).Delete(&Tags{}).Error
	if err != nil{
		return false, errors.New(helpers.LAN_DB_ERROR)
	}
	return true, nil
}
