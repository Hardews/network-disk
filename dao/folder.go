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
	err := dB.Create(&folder).Error
	if err != nil {
		return 0, err
	}
	err = dB.Create(&model.UserResources{
		FolderId:   uint(folder.ParentFolder),
		ResourceId: 0,
		Filename:   "folder",
		Permission: "folder",
	}).Error
	if err != nil {
		return 0, err
	}
	var folderId uint
	err = dB.Model(&model.Folder{}).Select("id").
		Where("username = ? and folder_name = ? and parent_folder = ?",
			folder.Username, folder.FolderName, folder.ParentFolder).Scan(&folderId).Error
	return folderId, err
}

func GetAllUserFolder(username string) []model.Folder {
	var res []model.Folder
	dB.Model(&model.Folder{}).Where("username = ?", username).Scan(&res)
	return res
}
