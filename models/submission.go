package models

import (
	"gorm.io/gorm"
)

type Submission struct {
	gorm.Model
	Identity        string  `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Path            string  `gorm:"column:path;type:varchar(255);" json:"path"`
	Status          int     `gorm:"column:status;type:tinyint;" json:"status"` // -1-待判断; 1-正确; 2-答案错误; 3-超时; 4-超内存; 5-编译错误
	UserIdentity    string  `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`
	User            User    `gorm:"foreignKey:user_identity;references:identity"`
	ProblemIdentity string  `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"`
	Problem         Problem `gorm:"foreignKey:problem_identity;references:identity"`
}

// GetSubmissionList 查询符合条件的提交记录
func GetSubmissionList(userIdentity string, problemIdentity string, status int) *gorm.DB {
	tx := DB.Model(&Submission{}).
		Preload("Problem", func(db *gorm.DB) *gorm.DB {
			return DB.Omit("content")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return DB.Omit("password")
		})
	if problemIdentity != "" {
		tx.Where("problem_identity = ?", problemIdentity)
	}
	if userIdentity != "" {
		tx.Where("user_identity = ?", userIdentity)
	}
	if status != 0 {
		tx.Where("status = ?", status)
	}
	return tx
}
