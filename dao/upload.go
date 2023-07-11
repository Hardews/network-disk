package dao

import (
	"network-disk/model"
	"time"
)

const (
	basePath        = "http://127.0.0.1:8080"
	oneDay          = 1 * 24 * time.Hour
	sevenDay        = 7 * oneDay
	foreverFileTime = 10 * 365 * 24 * time.Hour // 10 年
)

func ResourcesFile(ur model.UserResources) (bool, error) {
	// mysql
	err := dB.Create(&ur).Error
	if err != nil {
		return false, err
	}

	// redis
	key := ur.Path + "&&" + ur.Filename + "&&" + ur.Folder
	urStr := ur.ResourceName + "&&" + ur.Permission + "&&" + ur.CreateAt + "&&" + ur.DownloadAddr
	return rdb.HSet("user:"+ur.Username, key, urStr).Result()
}

func GetUrl(url string) (string, error) {
	res, err := rdb.Get(basePath + url).Result()
	if res == "" && err != nil {
		// redis 拿不到就从 mysql 拿
		var resUrl string
		err = dB.Model(&model.Url{}).Where("url = ? AND overdue >= ?", url, time.Now()).Scan(&resUrl).Error
		return resUrl, err
	}
	return res, err
}

func SetExpirationTime(url string, et int) error {
	var overdueTime time.Duration
	switch et {
	case 0:
		overdueTime = foreverFileTime
	case 1:
		rdb.Set(url, 1, oneDay)
		overdueTime = oneDay
	case 7:
		overdueTime = sevenDay
		rdb.Set(url, 7, sevenDay)
	}

	urlS := model.Url{Overdue: time.Now().Add(overdueTime), Url: url}
	tx := dB.Begin()

	dx := tx.Create(&urlS)
	if err := dx.Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func DelResourceFile(file model.UserResources) (int64, error) {
	dB.Model(&model.UserResources{}).Delete(&file)
	return rdb.HDel("user:"+file.Username, file.Path+"&&"+file.Filename+"&&"+file.Folder).Result()
}

func DbGetUserAllResource(username string) ([]model.UserResources, error) {
	var res []model.UserResources
	err := dB.Model(&model.UserResources{}).Where("username = ?", username).Scan(&res).Error
	return res, err
}

func DbGetUserResource(username, filename, path, folder string) (model.UserResources, error) {
	var res model.UserResources
	err := dB.Model(&model.UserResources{}).Where("username = ? AND filename = ? AND path = ? AND folder = ?",
		username, filename, path, folder).Scan(&res).Error
	return res, err
}

// redis

func RdbGetUserAllResource(username string) (map[string]string, error) {
	return rdb.HGetAll("user:" + username).Result()
}

func RdbGetUserResource(username, filename, path, folder string) (string, error) {
	return rdb.HGet("user:"+username, path+"&&"+filename+"&&"+folder).Result()
}

func ResourceIncr(resourceName string) (int64, error) {
	return rdb.Incr(resourceName).Result()
}

func ResourceDecr(resourceName string) (int64, error) {
	return rdb.Decr(resourceName).Result()
}
