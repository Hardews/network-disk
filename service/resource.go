/**
 * @Author: Hardews
 * @Date: 2023/7/14 15:50
 * @Description:
**/

package service

import (
	"network-disk/dao"
	"network-disk/model"
	"strconv"
	"strings"
)

func GetUserResourceByFolderId(folderId int) ([]model.UserResources, error) {
	return dao.GetUserResourceByFolderId(folderId)
}

func GetAllUserResource(username string) ([]model.UserResources, error) {
	urMap, err := dao.RdbGetUserAllResource(username)
	if err != nil {
		return nil, err
	}

	var urs []model.UserResources
	if urMap == nil {
		// redis 没获取到，去 MySQL 拿
		return dao.DbGetUserAllResource(username)
	} else {
		// 遍历得到的map，处理获取到的数据
		for key, ur := range urMap {
			first := strings.Split(key, ":")
			folderID, _ := strconv.Atoi(first[0])
			second := strings.Split(ur, "&&")
			resourceID, _ := strconv.Atoi(second[0])

			urs = append(urs, model.UserResources{
				FolderId:     uint(folderID),
				ResourceId:   uint(resourceID),
				Filename:     first[1],
				Permission:   second[1],
				DownloadAddr: second[2],
			})
		}
	}

	return urs, nil
}

// GetUserResource 根据信息来获取资源信息
func GetUserResource(filename string, folderId int) (ur model.UserResources, err error) {
	urStr, err := dao.RdbGetUserResource(filename, folderId)
	if err != nil {
		return
	}

	if urStr == "" {
		// 去 mysql 拿咯
		return dao.DbGetUserResource(filename, folderId)
	} else {
		// 解析 redis 中存储的字段
		second := strings.Split(urStr, "&&")
		resourceID, _ := strconv.Atoi(second[0])

		return model.UserResources{
			FolderId:     uint(folderId),
			ResourceId:   uint(resourceID),
			Filename:     filename,
			Permission:   second[1],
			DownloadAddr: second[2],
		}, nil
	}
}

func GetResourceName(resourceId int) string {
	return dao.GetResourceInfo(resourceId)
}

func GetResourceId(resourceName string) uint {
	return dao.GetResourceId(resourceName)
}
