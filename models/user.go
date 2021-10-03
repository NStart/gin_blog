package models

import (
	"fmt"
)

type AdminUser struct {
	BaseModel
	Name string `gorm:"type:varchar(20);unique;not null"`
	Password string `gorm:"type:varchar(32);"`
}

func (adminUser *AdminUser) GetOneUserByName(name string, oneUser *AdminUser) error {
	var err error
	db := GetDB()
	err = db.Where("name = ?", name).Take(&oneUser).Error
	fmt.Printf("%+v", oneUser)
	return err
}






