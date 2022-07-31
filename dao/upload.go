package dao

import "network-disk/model"

func ResourcesFile(username string, ur model.UserResources) (bool, error) {
	key := ur.Path + "&&" + ur.Filename + "&&" + ur.Folder
	urStr := ur.ResourceName + "&&" + ur.Permission + "&&" + ur.CreateAt + "&&" + ur.DownloadAddr
	return rdb.HSet("user:"+username, key, urStr).Result()
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
