package dao

import (
  	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

// 声明全局数据库变量DB
var DB *gorm.DB

func InitMySQL() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/projlist?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) 
	if err != nil {
		return
	}
	return 
}