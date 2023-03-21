package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

const (
	USERNAME = "root"
	PASSWORD = "alantam."
	HOST     = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "OJPlatform"
)

var DB = InitDB_()

func InitDB_() *gorm.DB {
	dsn := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8&parseTime=true"}, "")
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	//db.AutoMigrate(&TestCase{}, &Problem{})
	return db
}

/*
func PrepareSomeData() {
	var c Category
	c.Name = "array"
	c.Identity = "array"
	DB.Create(&c)

	var c2 Category
	c2.Name = "linkedList"
	c2.Identity = "linkedList"
	DB.Create(&c2)

	var p Problem
	p.Title = "prob 1"
	p.Content = "problem 111"
	p.Categories = append(p.Categories, c)
	DB.Create(&p)

	var p2 Problem
	p2.Title = "prob 2"
	p2.Content = "problem 222"
	p2.Categories = append(p2.Categories, c2)
	DB.Create(&p2)

	var p3 Problem
	p3.Title = "prob 3"
	p3.Content = "problem 333"
	p3.Categories = append(p3.Categories, c)
	DB.Create(&p3)
}
*/
