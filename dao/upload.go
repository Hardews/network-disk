package dao

import "network-disk/model"

func ResourcesFile(username string, ur model.UserResources) (bool, error) {
	urStr := ur.ResourceName + ":" + ur.Permission + ":" + ur.CreateAt + ":" + ur.Folder
	return rdb.HSet("user:"+username, ur.Filename, urStr).Result()
}

func DelResourceFile(username, filename string) (int64, error) {
	return rdb.HDel("user:"+username, filename).Result()
}

func GetUserAllResource(username string) (map[string]string, error) {
	return rdb.HGetAll("user:" + username).Result()
}

func GetUserResource(username, filename string) (string, error) {
	return rdb.HGet("user:"+username, filename).Result()
}

func ResourceIncr(resourceName string) (int64, error) {
	return rdb.Incr(resourceName).Result()
}

func ResourceDecr(resourceName string) (int64, error) {
	return rdb.Decr(resourceName).Result()
}
