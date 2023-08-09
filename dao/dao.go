package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"network-disk/config"
	"network-disk/model"
	"time"
)

var (
	dB *gorm.DB
)

func InitDB() {
	c := config.ReloadConfig
	var (
		mysqlUsername = c.MysqlDatabase.Username
		mysqlPassword = c.MysqlDatabase.Password
		mysqlLink     = c.DatabaseLink
		mysqlHost     = c.DatabaseHost
		mysqlName     = c.DatabaseName
	)

	// mysql link
	dsn := mysqlUsername + ":" + mysqlPassword + "@tcp(" + mysqlLink + ":" + mysqlHost + ")/" + mysqlName + "?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
	}
	dB = db

	dB.AutoMigrate(&model.Url{}, &model.UserResources{}, &model.Resource{}, &model.User{}, &model.Folder{}, &model.AdminUser{}, &model.Code{})

	// 创建文件夹的记录
	dB.Create(&model.Resource{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
		},
		ResourceName: "folder_link",
		ResourceNum:  0,
	})
}
