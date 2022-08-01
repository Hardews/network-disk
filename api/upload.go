package api

import (
	"encoding/base64"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"network-disk/model"
	"network-disk/service"
	"network-disk/tool"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/skip2/go-qrcode"
)

const basePath = "http://127.0.0.1:8080/"

func encryptionShare(ctx *gin.Context) {
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	pwd, res := ctx.GetQuery("password")
	var tmp string
	if !res {
		pwd = service.RandomStr(4)
		tmp = pwd
	}
	pwd = service.MD5([]byte(pwd))
	filename := ctx.Param("filename")

	storagePath, _ := ctx.GetQuery("path")
	var Path, err = base64.URLEncoding.DecodeString(storagePath)
	if err != nil {
		log.Println("encoding path failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	// 获取设定时间
	et, _ := ctx.GetQuery("time")
	eTime, err := strconv.Atoi(et)
	if err != nil {
		log.Println("upload:translate etime failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	folder, _ := ctx.GetQuery("category")

	str := pwd + "_" + filename + "_" + username + "_" + string(Path) + "_" + folder
	str = base64.URLEncoding.EncodeToString([]byte(str))

	url := basePath + "encryption/" + str
	ctx.JSON(http.StatusOK, gin.H{
		"password": tmp,
		"path":     url,
	})

	// 写入时间
	err = service.SetExpirationTime(url, eTime)
	if err != nil {
		log.Println("upload:set et failed,err:", err)
		tool.RespErrorWithDate(ctx, "分享失败,请重试")
		return
	}
}

// qrCode 二维码分享，用了github.com/skip2/go-qrcode包
func qrCode(ctx *gin.Context) {
	filename := ctx.Param("filename")

	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	folder, _ := ctx.GetQuery("category")
	storagePath, _ := ctx.GetQuery("path")
	var Path, err = base64.URLEncoding.DecodeString(storagePath)

	et, _ := ctx.GetQuery("time")
	eTime, err := strconv.Atoi(et)
	if err != nil {
		log.Println("upload:translate etime failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	ur, err := service.GetUserResource(username, filename, string(Path), folder)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	var (
		qr  *qrcode.QRCode
		url string
	)

	switch ur.Permission {
	case service.Public:
		url = basePath + username + "/" + filename + "?path=" + storagePath + "&category=" + folder
	case service.Private:
		tool.RespErrorWithDate(ctx, "您设置了仅自己可见，无法分享")
	case service.Permission:
		url = basePath + "download/"
		url += base64.URLEncoding.EncodeToString([]byte(filename + "-" + ur.ResourceName))
	}
	qr, err = qrcode.New(url, qrcode.Medium)
	if err != nil {
		log.Println(err)
		return
	}

	err = service.SetExpirationTime(url, eTime)
	if err != nil {
		log.Println("upload:set et failed,err:", err)
		tool.RespErrorWithDate(ctx, "分享失败,请重试")
		return
	}

	// 将二维码图生成并返回
	err = qr.Write(256, ctx.Writer)
}

func delFile(ctx *gin.Context) {
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	filename := ctx.PostForm("filename")

	Path := ctx.PostForm("path")

	folder := ctx.PostForm("category")

	ur, err := service.GetUserResource(username, filename, Path, folder)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	err = service.DelFile(username, ur.Filename, ur.ResourceName, Path, folder)
	if err != nil {
		if err == redis.Nil {
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
	// 获取各种信息
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	filename := ctx.PostForm("filename")
	Path := ctx.PostForm("path")
	folder := ctx.PostForm("category")
	// 这个是表明它想更改的是什么，路径，文件名。。。
	c := ctx.PostForm("chose")

	chose, err := strconv.Atoi(c)
	if err != nil {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	// 这个是应传入的新名称或路径
	newVal := ctx.PostForm("new")

	ur, err := service.GetUserResource(username, filename, Path, folder)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	res, err := service.UpdateFileAttribute(ur, newVal, username, chose)
	if err != nil {
		if err == service.ErrOfSameName {
			tool.RespErrorWithDate(ctx, err.Error())
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}
	if !res {
		tool.RespErrorWithDate(ctx, "更新失败")
		return
	}

	tool.RespSuccessful(ctx)
}

func downloadEncryptionFile(ctx *gin.Context) {
	pwd := ctx.PostForm("password")
	str, err := base64.URLEncoding.DecodeString(ctx.Param("filename"))

	s := strings.Split(string(str), "_")

	passwordSet := s[0]

	if service.MD5([]byte(pwd)) != passwordSet {
		tool.RespErrorWithDate(ctx, "密码错误")
		return
	}

	// 加密时url中存了这四个数据 username 是拥有者的id
	username := s[2]
	filename := s[1]
	Path := s[3]
	folder := s[4]

	ur, err := service.GetUserResource(username, filename, Path, folder)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	downloadFile(ctx, ur.Filename, ur.ResourceName)
}

func downloadPublicFile(ctx *gin.Context) {
	username := ctx.Param("username")
	filename := ctx.Param("filename")

	storagePath, _ := ctx.GetQuery("path")
	var Path, err = base64.URLEncoding.DecodeString(storagePath)
	folder, _ := ctx.GetQuery("category")

	ur, err := service.GetUserResource(username, filename, string(Path), folder)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	if ur.Permission != service.Public {
		ctx.JSON(http.StatusForbidden, "没有权限下载该文件")
		return
	}

	downloadFile(ctx, filename, ur.ResourceName)
}

func downloadFileByConn(ctx *gin.Context) {

	name := ctx.Param("filename")

	str, err := base64.URLEncoding.DecodeString(name)
	if err != nil {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	s := strings.Split(string(str), "-")
	filename := s[0]
	resource := s[1]

	downloadFile(ctx, filename, resource)
}

// downloadUserFile 用户下载自己的资源
func downloadUserFile(ctx *gin.Context) {
	// 获取各种东西
	filename := ctx.Param("filename")

	storagePath, _ := ctx.GetQuery("path")
	var Path, err = base64.URLEncoding.DecodeString(storagePath)

	folder, _ := ctx.GetQuery("category")

	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	ur, err := service.GetUserResource(username, filename, string(Path), folder)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	downloadFile(ctx, filename, ur.ResourceName)
}

// downloadFile 文件下载
func downloadFile(ctx *gin.Context, filename, resource string) {
	file, err := os.Open(resource)
	if err != nil {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	_, err = file.Read(fileHeader)
	if err != nil {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	fileStat, err := file.Stat()
	if err != nil {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	// 设置响应头
	ctx.Writer.Header().Set("Content-Disposition", "attachment;filename="+filename)
	ctx.Writer.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	ctx.Writer.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	// 设置初始偏移量
	file.Seek(0, 0)

	for {
		var n int
		// 通过控制切片大小控制下载速度
		tmp := make([]byte, 125*100) // 1 兆
		n, err = file.Read(tmp)
		if err == io.EOF {
			return
		}

		ctx.Writer.Write(tmp[:n])
		time.Sleep(1 * time.Millisecond)
	}
}

// shareFile 分享文件
func shareFile(ctx *gin.Context) {
	// 获取用户名
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)
	// 获取文件名
	filename := ctx.Param("filename")
	// 获取路径
	storagePath, _ := ctx.GetQuery("path")
	var Path, err = base64.URLEncoding.DecodeString(storagePath)
	// 获取文件夹位置
	folder, _ := ctx.GetQuery("category")
	// 获取过期时间
	et, _ := ctx.GetQuery("time")
	expirationTime, err := strconv.Atoi(et)
	if err != nil {
		log.Println("upload:translate et failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	// 获取该用户想分享的资源的信息
	ur, err := service.GetUserResource(username, filename, string(Path), folder)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	var (
		str string
		url string
	)

	switch ur.Permission {
	case service.Public:
		url = basePath + username + "/" + filename + "?path=" + storagePath + "&category=" + folder
	case service.Private:
		tool.RespErrorWithDate(ctx, "分享失败，您以将该文件设置为仅自己可见")
	case service.Permission:
		// 使人们只能通过分享连接下载的想法是将url进行base64编码
		str = base64.URLEncoding.EncodeToString([]byte(filename + "-" + ur.ResourceName))
		url = basePath + "download/" + str
	}
	err = service.SetExpirationTime(url, expirationTime)
	if err != nil {
		log.Println("upload:set et failed,err:", err)
	}
	tool.RespSuccessfulWithDate(ctx, url)
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

	// 获取存储的文件夹
	folder := ctx.PostForm("category")
	if folder == "" {
		tool.RespErrorWithDate(ctx, "未指定文件夹")
		return
	}

	// 获取存储的路径
	Path := ctx.PostForm("path")
	if Path == "" {
		tool.RespErrorWithDate(ctx, "未指定路径")
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
	res, err = service.IsRepeatFilename(username, file.Filename, folder, Path)
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

	var storagePath = base64.URLEncoding.EncodeToString([]byte(Path))
	// 在redis中存储的结构
	var storage = model.UserResources{
		Folder:       folder,
		Path:         Path,
		Filename:     filename,
		ResourceName: resourceName,
		Permission:   attribute,
		DownloadAddr: basePath + "user/download/" + file.Filename + "?path=" + storagePath + "&category=" + folder,
		CreateAt:     time.Now().String(),
	}

	// 存储在redis中
	_, err = service.StorageFile(username, storage)
	if err != nil {
		log.Println("storage file failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
}

// getUserFileByCategory 获取指定路径的文件
func getUserFileByCategory(ctx *gin.Context) {
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	// 获取文件夹名称
	category, res := ctx.GetQuery("category")
	if !res {
		tool.RespErrorWithDate(ctx, "文件夹名称为空")
		return
	}

	// 获取路径
	Path, res := ctx.GetQuery("path")
	if !res {
		tool.RespErrorWithDate(ctx, "路径为空")
		return
	}

	// 获取想要的文件信息
	urs, err := service.GetUserFileByCategory(username, category, Path)
	if err != nil {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}
	tool.RespSuccessfulWithDate(ctx, urs)
}

// getUserAllFile 获取用户的所有文件
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
