/**
 * @Author: Hardews
 * @Date: 2023/7/14 15:50
 * @Description:
**/

package service

import (
	"log"
	"network-disk/dao"
	"network-disk/model"
	"os"
	"strconv"
	"strings"
)

// UpdateFileAttribute 更新文件属性
//func UpdateFileAttribute(old model.UserResources, new string, chose int) (res bool, err error) {
//
//}

// DelFile 删除文件
func DelFile(folderId int, filename string) (err error) {
	var ResourceId int
	// 检查该用户是否有存储该文件
	resStr, err := dao.RdbGetUserResource(filename, folderId)
	if err != nil {
		log.Println("redis get user resource failed,err:", err)
		// 去 mysql 拿
		res, err := dao.DbGetUserResource(filename, folderId)
		if res.Filename == "" || err != nil {
			// mysql 也没拿到
			return err
		}
		ResourceId = int(res.ResourceId)
	} else {
		resArr := strings.Split(resStr, "&&")
		ResourceId, err = strconv.Atoi(resArr[0])
		if err != nil {
			return err
		}
	}

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

	_, err = dao.DelResourceFile(model.UserResources{
		FolderId:   uint(folderId),
		ResourceId: uint(ResourceId),
		Filename:   filename,
	})
	return
}
