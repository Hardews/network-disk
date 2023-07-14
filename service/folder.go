/**
 * @Author: Hardews
 * @Date: 2023/7/14 17:32
 * @Description:
**/

package service

import (
	"network-disk/dao"
	"network-disk/model"
)

func CreateFolder(folder model.Folder) (uint, error) {
	return dao.CreateFolder(folder)
}

func GetAllFolder(username string) []model.Folder {
	return dao.GetAllUserFolder(username)
}
