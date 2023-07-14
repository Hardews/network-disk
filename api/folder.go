/**
 * @Author: Hardews
 * @Date: 2023/7/14 17:31
 * @Description:
**/

package api

import (
	"github.com/gin-gonic/gin"

	"log"
	"strconv"

	"network-disk/model"
	"network-disk/service"
	"network-disk/tool"
)

func getFolderInfo(ctx *gin.Context) {
	iUsername, exist := ctx.Get("username")
	if !exist {
		tool.RespSuccessfulWithDate(ctx, "登陆后再操作")
		return
	}
	username := iUsername.(string)
	tool.RespSuccessfulWithDate(ctx, service.GetAllFolder(username))
}

func addFolder(ctx *gin.Context) {
	iUsername, exist := ctx.Get("username")
	if !exist {
		tool.RespSuccessfulWithDate(ctx, "登陆后再操作")
		return
	}
	username := iUsername.(string)

	parentFolder, exists := ctx.GetPostForm("parent")
	if !exists {
		tool.RespSuccessfulWithDate(ctx, "父文件夹 id 为空")
		return
	}
	parentId, err := strconv.Atoi(parentFolder)
	if err != nil {
		tool.RespSuccessfulWithDate(ctx, "传入参数错误")
		return
	}

	folderName, exists := ctx.GetPostForm("folder_name")
	if !exists {
		tool.RespSuccessfulWithDate(ctx, "文件夹名字为空")
		return
	}

	id, err := service.CreateFolder(model.Folder{
		Username:     username,
		FolderName:   folderName,
		ParentFolder: uint(parentId),
	})
	if err != nil {
		log.Println("create folder err:", err)
		tool.RespInternetError(ctx)
	}
	tool.RespSuccessfulWithDate(ctx, id)
}
