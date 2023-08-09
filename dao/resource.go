/**
 * @Author: Hardews
 * @Date: 2023/7/14 17:23
 * @Description:
**/

package dao

import (
	"network-disk/model"
)

func CreateANewResource(resourceName string) error {
	return dB.Create(&model.Resource{
		ResourceName: resourceName,
		ResourceNum:  1,
	}).Error
}

func ResourcesFile(ur model.UserResources) (bool, error) {
	err := dB.Create(&ur).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func DelResourceFile(file model.UserResources) error {
	return dB.Model(&model.UserResources{}).Delete(&file).Error
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
	return nowNum, nil
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

	return nowNum, nil
}
