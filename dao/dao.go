package dao

import (
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	dB  *gorm.DB
	rdb *redis.Client
)

func InitDB() {
	dsn := "root:lmh123@tcp(127.0.0.1:3306)/disk?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
	}
	dB = db

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "lmh123",
		DB:       10,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		log.Fatalln(err)
	}
}
