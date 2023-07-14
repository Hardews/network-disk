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
	res, err := rdb.Get(basePath + url).Result()
	if res == "" && err != nil {
		// redis 拿不到(不存在或过期）就从 mysql 拿
		var resUrl string
		err = dB.Model(&model.Url{}).Where("url = ? AND overdue >= ?", url, time.Now()).Scan(&resUrl).Error
		return resUrl, err
	}
	return res, err
}

func SetExpirationTime(url string, et int, code ...string) error {
	var overdueTime time.Duration
	switch et {
	case 0:
		overdueTime = foreverFileTime
	case 1:
		if len(code) != 0 {
			rdb.Set(url, code[0], oneDay)
		} else {
			rdb.Set(url, 1, oneDay)
		}
		overdueTime = oneDay
	case 7:
		if len(code) != 0 {
			rdb.Set(url, code[0], sevenDay)
		} else {
			rdb.Set(url, 7, sevenDay)
		}
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
