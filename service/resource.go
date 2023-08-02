/**
 * @Author: Hardews
 * @Date: 2023/7/14 15:50
 * @Description:
**/

package service

import (
	"network-disk/dao"
	"network-disk/model"
)

func GetUserResourceByFolderId(folderId int) ([]model.UserResources, error) {
	return dao.GetUserResourceByFolderId(folderId)
}

func GetAllUserResource(username string) ([]model.UserResources, error) {
	// redis 没获取到，去 MySQL 拿
	return dao.DbGetUserAllResource(username)
}

// GetUserResource 根据信息来获取资源信息
func GetUserResource(filename string, folderId int) (ur model.UserResources, err error) {
	// 去 mysql 拿咯
	return dao.DbGetUserResource(filename, folderId)
}

func GetResourceName(resourceId int) string {
	return dao.GetResourceInfo(resourceId)
}

func GetResourceId(resourceName string) uint {
	return dao.GetResourceId(resourceName)
}
