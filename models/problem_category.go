package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	ProblemId  uint     `gorm:"column:problem_id;" json:"problem_id"`
	CategoryId uint     `gorm:"column:category_id;" json:"category_id"`
	Category   Category `gorm:"foreignKey:category_id;references:id"`
}
