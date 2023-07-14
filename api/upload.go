package api

import (
	"errors"
	"gorm.io/gorm"
	"io/fs"
	"log"
	"network-disk/config"
	"os"
	"path"
	"strconv"
	"time"

	"network-disk/model"
	"network-disk/service"
	"network-disk/tool"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var basePath = config.ReloadConfig.BaseUrl + "/"

func delFile(ctx *gin.Context) {
	filename := ctx.PostForm("filename")

	iFolderId := ctx.PostForm("folder_id")

	folderId, err := strconv.Atoi(iFolderId)
	if err != nil {
		tool.RespErrorWithDate(ctx, "param error")
		return
	}

	err = service.DelFile(folderId, filename)
	if err != nil {
		if errors.Is(err, redis.Nil) || errors.Is(err, gorm.ErrRecordNotFound) {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
}

// updateFileAttribute 更新文件的属性
func updateFileAttribute(ctx *gin.Context) {
	//res, err := service.UpdateFileAttribute(ur, newVal, chose)
	//if err != nil {
	//	if err == service.ErrOfSameName {
	//		tool.RespErrorWithDate(ctx, err.Error())
	//		return
	//	}
	//	log.Println(err)
	//	tool.RespInternetError(ctx)
	//	return
	//}
	//if !res {
	//	tool.RespErrorWithDate(ctx, "更新失败")
	//	return
	//}
	//
	//tool.RespSuccessful(ctx)
}

// uploadFile 上传文件
func uploadFile(ctx *gin.Context) {
	// 获取用户名
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	// 获取属性,不设置则为公开
	attribute := ctx.PostForm("attribute")
	if attribute == "" {
		attribute = service.Public
	}

	// 获取存储的文件夹 id
	iFolderId := ctx.PostForm("folder_id")
	folderId, err := strconv.Atoi(iFolderId)
	if err != nil {
		tool.RespErrorWithDate(ctx, "传参错误")
		return
	}
	if folderId <= 0 {
		tool.RespErrorWithDate(ctx, "未指定上传文件夹")
		return
	}

	// 获取要上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Println("get file failed,err:", err)
		tool.RespErrorWithDate(ctx, "上传失败，请重试")
		return
	}
	if file == nil {
		tool.RespErrorWithDate(ctx, "文件为空")
		return
	}

	// 预处理文件
	res, resourceName, err := service.DealWithFile(file)
	if !res {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}
	if err != nil {
		tool.RespSuccessfulWithDate(ctx, err)
		return
	}

	// 判断是否与已存在资源的内容重复
	// 重复: 只存一个对应的连接
	// 不重复: 存在服务器

	var (
		filename       = file.Filename
		breakFile      *os.File
		IsFileRepeat   = true
		resourceInfo   fs.FileInfo
		fileSuffix     = path.Ext(file.Filename)
		breakPointPath = "./uploadFile/breakPoint/" + username + filename[:len(filename)-len(fileSuffix)] + ".txt"
	)

	resource, err := os.Open(resourceName)
	if err != nil {
		if os.IsNotExist(err) {
			IsFileRepeat = false
			err = nil
			goto storage
		} else {
			log.Println("open resource file failed,err:", err)
			tool.RespErrorWithDate(ctx, "上传失败，请重试")
			return
		}
	}

	resourceInfo, err = resource.Stat()
	if err != nil {
		log.Println("get resource info failed,err:", err)
		tool.RespErrorWithDate(ctx, "上传失败，请重试")
		return
	}

storage:
	// 如果有这个文件但是这个文件的大小与上传的不一致证明中断过
	if !IsFileRepeat || resourceInfo.Size() != file.Size {
		err = service.Storage(file, resourceName, breakPointPath, breakFile)
		if err != nil {
			log.Println(err)
			tool.RespInternetError(ctx)
			return
		}
	}

	// 这里是判断用户上传的文件名是否与存在的文件名重复
	// 如果重复，帮它改名字
	res, err = service.IsRepeatFilename(file.Filename, folderId)
	if err != nil {
		log.Println("判断重复名失败,err:", err)
		tool.RespErrorWithDate(ctx, "服务器错误，上传文件失败")
		return
	}
	if !res {
		fileSuffix := path.Ext(file.Filename)
		Len := len(file.Filename) - len(fileSuffix)
		filename = file.Filename[:Len] + time.Now().Format("20060102_030405") + fileSuffix
	}

	// 获取对应 resource 的 id

	// 在数据库中存储的结构
	var storage = model.UserResources{
		FolderId:     uint(folderId),
		ResourceId:   service.GetResourceId(resourceName),
		Filename:     filename,
		Permission:   attribute,
		DownloadAddr: basePath + "posts/download/" + strconv.Itoa(folderId) + "/?filename=" + filename,
	}

	// 存储在 redis, mysql 中
	_, err = service.StorageFile(storage)
	if err != nil {
		log.Println("storage file failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
}

// getUserFileByCategory 获取指定路径的文件
func getUserFileByCategory(ctx *gin.Context) {
	// 获取存储的文件夹 id
	iFolderId := ctx.PostForm("folder_id")
	folderId, err := strconv.Atoi(iFolderId)
	if err != nil {
		tool.RespErrorWithDate(ctx, "传参错误")
		return
	}
	if folderId <= 0 {
		tool.RespErrorWithDate(ctx, "未指定上传文件夹")
		return
	}

	// 获取想要的文件信息
	urs, err := service.GetUserResourceByFolderId(folderId)
	if err != nil {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}
	// 返回文件夹下的所有文件信息
	tool.RespSuccessfulWithDate(ctx, urs)
}

// getUserAllFile 获取所有文件
func getUserAllFile(ctx *gin.Context) {
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	urs, err := service.GetAllUserResource(username)
	if err != nil {
		log.Println(err)
		return
	}

	tool.RespSuccessfulWithDate(ctx, urs)
}
