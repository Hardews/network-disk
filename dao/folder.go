/**
 * @Author: Hardews
 * @Date: 2023/7/14 17:24
 * @Description:
**/

package dao

import (
	"network-disk/model"
)

func GetUsernameByFolderId(folderId uint) string {
	// 获取用户名
	var username string
	dB.Model(&model.Folder{}).Select("username").Where("id = ?", folderId).Scan(&username)
	return username
}

func CreateFolder(folder model.Folder) (uint, error) {
	// 需要看看这个文件夹名字有没有出现
	var isFolderRepeat string
	dB.Model(&model.Folder{}).Select("folder_name").Where("username = ? and folder_name = ?",
		folder.Username, folder.FolderName).Scan(&isFolderRepeat)
	// 一致则加个副本字段
	if isFolderRepeat == folder.FolderName {
		folder.FolderName = folder.FolderName + "_副本"
	}

	// Mysql 插入数据
	err := dB.Create(&folder).Error
	if err != nil {
		return 0, err
	}

	// mysql 更新
	err = dB.Create(&model.UserResources{
		FolderId:   uint(folder.ParentFolder),
		ResourceId: 0,
		Filename:   "folder",
		Permission: "folder",
	}).Error
	if err != nil {
		return 0, err
	}

	// 查询 folder id
	var folderId uint
	err = dB.Model(&model.Folder{}).Select("id").
		Where("username = ? and folder_name = ? and parent_folder = ?",
			folder.Username, folder.FolderName, folder.ParentFolder).Scan(&folderId).Error
	if err != nil {
		return 0, err
	}

	return folderId, err
}

func GetAllUserFolder(username string) []model.Folder {
	var res []model.Folder
	dB.Model(&model.Folder{}).Where("username = ?", username).Scan(&res)
	return res
}
