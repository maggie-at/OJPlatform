package models

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Identity  string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Name      string `gorm:"column:name;type:varchar(100);" json:"name"`
	Password  string `gorm:"column:password;type:varchar(32);" json:"password"`
	Phone     string `gorm:"column:phone;type:varchar(100);" json:"phone"`
	Email     string `gorm:"column:email;type:varchar(100);" json:"email"`
	PassNum   int64  `gorm:"column:pass_num;type:int(11);" json:"pass_num"`     // 总通过次数
	SubmitNum int64  `gorm:"column:submit_num;type:int(11);" json:"submit_num"` // 总提交次数
	IsAdmin   int    `gorm:"column:is_admin;type:tinyint(1);" json:"is_admin"`  // 0-普通用户; 1-提交次数
}

func GetUserDetail(identity string) (*User, error) {
	var u User
	// 省略密码字段
	tx := DB.Omit("password").Where("identity = ?", identity).First(&u)
	if err := tx.Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func UserLogin(username string, password string) (*User, error) {
	var u User
	// 不返回敏感信息
	tx := DB.Select("name", "identity", "is_admin").Where("name = ? AND password = ?", username, password).First(&u)
	if err := tx.Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func CheckEmailExist(email string) (int64, error) {
	var cnt int64
	err := DB.Model(&User{}).Where("email = ?", email).Count(&cnt).Error
	fmt.Println(cnt)
	return cnt, err
}
func UserRegister(user User) error {
	return DB.Create(&user).Error
}
func GeuUserRank() *gorm.DB {
	tx := DB.Model(&User{}).Select("name", "identity", "finish_problem_num", "submit_num").Order("finish_problem_num DESC, submit_num ASC")
	return tx
}
