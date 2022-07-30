package dao

import "network-disk/model"

func ResourcesFile(username string, ur model.UserResources) (bool, error) {
	urStr := ur.ResourceName + ":" + ur.Permission + ":" + ur.CreateAt
	return rdb.HSet("user:"+username, ur.Filename, urStr).Result()
}

func GetUserAllResource(username string) (map[string]string, error) {
	return rdb.HGetAll("user:" + username).Result()
}

func GetUserResource(username, filename string) (string, error) {
	return rdb.HGet("user:"+username, filename).Result()
}
