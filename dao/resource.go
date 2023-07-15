/**
 * @Author: Hardews
 * @Date: 2023/7/14 17:23
 * @Description:
**/

package dao

import (
	"network-disk/model"
	"strconv"
)

func CreateANewResource(resourceName string) error {
	return dB.Create(&model.Resource{
		ResourceName: resourceName,
		ResourceNum:  0,
	}).Error
}

func ResourcesFile(ur model.UserResources) (bool, error) {
	// mysql
	err := dB.Create(&ur).Error
	if err != nil {
		return false, err
	}

	username := GetUsernameByFolderId(ur.FolderId)

	/*
	 redis
	 哈希组名称 ： user:userId
	 key : folderID + 文件名称
	 val ：resourceId + 权限 + 下载地址
	*/
	key := strconv.Itoa(int(ur.FolderId)) + ":" + ur.Filename
	urStr := strconv.Itoa(int(ur.ResourceId)) + "&&" + ur.Permission + "&&" + ur.DownloadAddr
	return rdb.HSet(redisStoragePrefix+username, key, urStr).Result()
}

func DelResourceFile(file model.UserResources) (int64, error) {
	dB.Model(&model.UserResources{}).Delete(&file)
	hashKey := redisStoragePrefix + GetUsernameByFolderId(file.FolderId)
	key := strconv.Itoa(int(file.FolderId)) + file.Filename
	return rdb.HDel(hashKey, key).Result()
}

func GetResourceId(resourceName string) uint {
	var resourceId uint
	dB.Model(&model.Resource{}).Select("id").Where("resource_name = ?", resourceName).Scan(&resourceId)
	return resourceId
}

func DbGetUserAllResource(username string) ([]model.UserResources, error) {
	var ids []uint
	err := dB.Model(&model.Folder{}).Select("id").Where("username = ?", username).Scan(&ids).Error
	if err != nil {
		return nil, err
	}

	var res []model.UserResources
	for _, id := range ids {
		var temp []model.UserResources
		err = dB.Model(&model.UserResources{}).Where("folder_id = ?", id).Scan(&temp).Error
		if err != nil {
			return nil, err
		}
		if temp != nil {
			res = append(res, temp...)
		}
	}

	return res, err
}

func GetUserResourceByFolderId(folderId int) ([]model.UserResources, error) {
	var res []model.UserResources
	err := dB.Model(&model.UserResources{}).Where("folder_id = ?", folderId).Scan(&res).Error
	return res, err
}

func DbGetUserResource(filename string, folderId int) (model.UserResources, error) {
	var res model.UserResources
	err := dB.Model(&model.UserResources{}).Where("folder_id = ? AND filename = ?",
		folderId, filename).Scan(&res).Error
	return res, err
}

func GetResourceInfo(resourceId int) string {
	var res string
	dB.Model(&model.Resource{}).Select("resource_name").Where("id = ?", resourceId).Scan(&res)
	return res
}

// redis

func RdbGetUserAllResource(username string) (map[string]string, error) {
	return rdb.HGetAll(redisStoragePrefix + username).Result()
}

func RdbGetUserResource(filename string, folderId int) (string, error) {
	key := redisStoragePrefix + GetUsernameByFolderId(uint(folderId))
	val := strconv.Itoa(folderId) + ":" + filename
	return rdb.HGet(key, val).Result()
}

// ResourceIncr 记录资源数++，当没人拥有这个文件时删除
func ResourceIncr(resourceId uint) (int64, error) {
	var nowNum int64
	dB.Model(&model.Resource{}).Select("resource_num").Where("id = ?", resourceId).Scan(&nowNum)
	if nowNum != 0 {
		nowNum++
		err := dB.Model(&model.Resource{}).Where("id = ?", resourceId).Update("resource_num", nowNum).Error
		if err != nil {
			return -1, err
		}
	}
	// 更新 redis 并处理数据一致性的问题
	resNum, err := rdb.Incr(strconv.Itoa(int(resourceId))).Result()
	if err != nil {
		return -1, err
	}

	if resNum != nowNum {
		// 以 mysql 的为准
		rdb.Set(strconv.Itoa(int(resourceId)), nowNum, -1)
	}
	return nowNum, err
}

func ResourceDecr(resourceId uint) (int64, error) {
	var nowNum int64
	dB.Model(&model.Resource{}).Select("resource_num").Where("id = ?", resourceId).Scan(&nowNum)
	if nowNum != 0 {
		nowNum--
		err := dB.Model(&model.Resource{}).Where("id = ?", resourceId).Update("resource_num", nowNum).Error
		if err != nil {
			return -1, err
		}
	}
	// 更新 redis 并处理数据一致性的问题
	resNum, err := rdb.Decr(strconv.Itoa(int(resourceId))).Result()
	if err != nil {
		return -1, err
	}

	if resNum != nowNum {
		// 以 mysql 的为准
		rdb.Set(strconv.Itoa(int(resourceId)), nowNum, -1)
	}
	return nowNum, err
}
