package dao

import (
	"network-disk/model"
	"time"
)

const (
	basePath           = "http://127.0.0.1:8080"
	oneDay             = 1 * 24 * time.Hour
	sevenDay           = 7 * oneDay
	foreverFileTime    = 10 * 365 * 24 * time.Hour // 10 年
	redisStoragePrefix = "user:"                   // redis 哈希组存储的前缀
)

func GetUrl(url string) (string, error) {
	var res string
	err := dB.Model(&model.Url{}).Where("url = ? AND overdue >= ?", url, time.Now()).Scan(&res).Error
	return res, err
}

func SetExpirationTime(url string, et int, code ...string) error {
	var overdueTime time.Duration
	switch et {
	case 0:
		overdueTime = foreverFileTime
	case 1:
		overdueTime = oneDay
	case 7:
		overdueTime = sevenDay
	}

	urlS := model.Url{Overdue: time.Now().Add(overdueTime), Url: url}
	tx := dB.Begin()
	defer tx.Rollback()

	dx := tx.Create(&urlS)
	if err := dx.Error; err != nil {
		return err
	}

	if len(code) != 0 {
		tx.Create(&model.Code{
			Url:  url,
			Code: code[0],
		})
	}

	tx.Commit()
	return nil
}
