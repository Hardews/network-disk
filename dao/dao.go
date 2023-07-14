package dao

import (
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"network-disk/config"
	"network-disk/model"
)

var (
	dB  *gorm.DB
	rdb *redis.Client
)

func InitDB() {
	c := config.ReloadConfig
	var (
		mysqlUsername = c.MysqlDatabase.Username
		mysqlPassword = c.MysqlDatabase.Password
		mysqlLink     = c.DatabaseLink
		mysqlHost     = c.DatabaseHost
		mysqlName     = c.DatabaseName

		redisAddr     = c.Address
		redisPassword = c.RedisDatabase.Password
		redisDb       = c.RedisDatabase.DB
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

	// redis link
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDb,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		log.Println("failed to connect redis,err:", err)
		// log.Fatalln("failed to connect redis,err:", err)
	}
}
