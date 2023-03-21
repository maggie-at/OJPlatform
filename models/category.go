package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Name     string `gorm:"column:name;type:varchar(100);" json:"name"`
	ParentId int    `gorm:"column:parent_id;type:int(11);" json:"parent_id"`
	LinkUrl  string `gorm:"column:link_url;type:varchar(255);" json:"link_url"`
}

func GetCategoryList(keyword string) *gorm.DB {
	if keyword != "" {
		return DB.Model(&Category{}).Where("name like ?", "%"+keyword+"%")
	} else {
		return DB.Model(&Category{}).Order("name ASC")
	}
}
