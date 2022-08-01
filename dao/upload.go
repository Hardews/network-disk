package dao

import (
	"network-disk/model"
	"time"
)

func ResourcesFile(username string, ur model.UserResources) (bool, error) {
	key := ur.Path + "&&" + ur.Filename + "&&" + ur.Folder
	urStr := ur.ResourceName + "&&" + ur.Permission + "&&" + ur.CreateAt + "&&" + ur.DownloadAddr
	return rdb.HSet("user:"+username, key, urStr).Result()
}

func GetUrl(url string) (string, error) {
	return rdb.Get(url).Result()
}

func GetForeverUrl(url string) error {
	tx := dB.Where("url = ?", url).First(&model.Url{})
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func SetExpirationTime(url string, et int) error {
	switch et {
	case 0:
		urlS := model.Url{Url: url}
		tx := dB.Begin()

		dx := tx.Create(&urlS)
		if err := dx.Error; err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
	case 1:
		rdb.Set(url, 1, 1*24*time.Hour)
	case 7:
		rdb.Set(url, 7, 7*24*time.Hour)
	}
	return nil
}

func DelResourceFile(username, filename, path, folder string) (int64, error) {
	return rdb.HDel("user:"+username, path+"&&"+filename+"&&"+folder).Result()
}

func GetUserAllResource(username string) (map[string]string, error) {
	return rdb.HGetAll("user:" + username).Result()
}

func GetUserResource(username, filename, path, folder string) (string, error) {
	return rdb.HGet("user:"+username, path+"&&"+filename+"&&"+folder).Result()
}

func ResourceIncr(resourceName string) (int64, error) {
	return rdb.Incr(resourceName).Result()
}

func ResourceDecr(resourceName string) (int64, error) {
	return rdb.Decr(resourceName).Result()
}
