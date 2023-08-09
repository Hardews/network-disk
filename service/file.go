/**
 * @Author: Hardews
 * @Date: 2023/7/14 15:50
 * @Description:
**/

package service

import (
	"network-disk/dao"
	"network-disk/model"
	"os"
)

// UpdateFileAttribute 更新文件属性
//func UpdateFileAttribute(old model.UserResources, new string, chose int) (res bool, err error) {
//
//}

// DelFile 删除文件
func DelFile(folderId int, filename string) (err error) {
	var ResourceId int
	res, err := dao.DbGetUserResource(filename, folderId)
	if res.Filename == "" || err != nil {
		// mysql 也没拿到
		return err
	}
	ResourceId = int(res.ResourceId)

	n, err := dao.ResourceDecr(uint(ResourceId))
	if err != nil {
		return
	}
	if n == 0 {
		// 没有人存储这个文件了，删除
		ResourceName := dao.GetResourceInfo(ResourceId)
		err = os.Remove(ResourceName)
		if err != nil {
			return
		}
	}
	return dao.DelResourceFile(model.UserResources{
		FolderId:   uint(folderId),
		ResourceId: uint(ResourceId),
		Filename:   filename,
	})
}
