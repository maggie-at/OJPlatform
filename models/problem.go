package models

import (
	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model
	Identity          string            `gorm:"column:identity;type:varchar(36);" json:"identity"`             // 问题的唯一标识
	Title             string            `gorm:"column:title;type:varchar(255);" json:"title"`                  // 题目标题
	Content           string            `gorm:"column:content;type:text;" json:"content"`                      // 题目描述
	MaxMemory         int               `gorm:"colum:max_memory;type:int;" json:"max_memory"`                  // 最大运行内存(KB)
	MaxRuntime        int               `gorm:"colum:max_runtime;type:int;" json:"max_runtime"`                // 最大运行时间(ms)
	PassNum           int64             `gorm:"column:pass_num;type:int(11);" json:"pass_num"`                 // 通过次数
	SubmitNum         int64             `gorm:"column:submit_num;type:int(11);" json:"submit_num"`             // 提交次数
	ProblemCategories []ProblemCategory `gorm:"foreignKey:problem_id;references:id" json:"problem_categories"` // 多对多关系
	TestCases         []TestCase        `gorm:"foreignKey:problem_identity;references:identity"`               // 一对多关系
}

// GetProblemList 根据关键词模糊查询
func GetProblemList(keyword string, categoryIdentity string) *gorm.DB {
	tx := DB.Model(&Problem{}).Preload("ProblemCategories").
		Where("title like ? or content like ?", "%"+keyword+"%", "%"+keyword+"%")
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_categories as PC on PC.problem_id = problems.id").
			Where("PC.category_id = ( select C.id from categories as C where C.identity = ? )", categoryIdentity)
	}
	return tx
}

func GetProblemDetail(identity string) (*Problem, error) {
	var p Problem
	// 查询Categories列表, 需要.Preload("字段名"), 需要和结构体中的字段名严格一致 https://segmentfault.com/a/1190000043248394
	tx := DB.Where("identity = ?", identity).Preload("Categories").First(&p) // 联合查询(为了查出来categories切片)
	if err := tx.Error; err != nil {
		return nil, err
	}
	return &p, nil
}
