package test

import (
	"OJPlatform/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"testing"
)

const (
	USERNAME = "root"
	PASSWORD = "alantam."
	HOST     = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "OJPlatform"
)

func TestGormTest(t *testing.T) {
	dsn := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8&parseTime=true"}, "")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Problem{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Submission{})
}
