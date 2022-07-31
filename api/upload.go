package api

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
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
	pwd := ctx.Request.Header.Get("password")
	var tmp string
	if pwd == "" {
		pwd = service.RandomStr(4)
		tmp = pwd
	}
	pwd = service.MD5([]byte(pwd))
	filename := ctx.Param("filename")

	str := pwd + "_" + filename
	str = base64.URLEncoding.EncodeToString([]byte(str))

	ctx.JSON(http.StatusOK, gin.H{
		"password": tmp,
		"path":     basePath + "encryption/" + str,
	})
}

func qrCode(ctx *gin.Context) {
	filename := ctx.Param("filename")
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	ur, err := service.GetUserResource(username, filename)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	var qr *qrcode.QRCode
	switch ur.Permission {
	case service.Public:
		qr, err = qrcode.New(basePath+username+"/"+filename, qrcode.Medium)
	case service.Private:
		tool.RespErrorWithDate(ctx, "您设置了仅自己可见，无法分享")
	case service.Permission:
		str := basePath + "download/"
		str += base64.URLEncoding.EncodeToString([]byte(filename + "-" + ur.ResourceName))
		qr, err = qrcode.New(str, qrcode.Medium)
	}
	if err != nil {
		log.Println(err)
		return
	}

	err = qr.Write(256, ctx.Writer)
}

func delFile(ctx *gin.Context) {
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)
	filename := ctx.PostForm("filename")

	ur, err := service.GetUserResource(username, filename)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	err = service.DelFile(username, ur.Filename, ur.ResourceName)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "您没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
}

func updateFileAttribute(ctx *gin.Context) {
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)
	filename := ctx.PostForm("filename")
	c := ctx.PostForm("chose")

	chose, err := strconv.Atoi(c)
	if err != nil {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	newVal := ctx.PostForm("new")

	ur, err := service.GetUserResource(username, filename)
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
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	s := strings.Split(string(str), "_")

	if service.MD5([]byte(pwd)) != s[0] {
		tool.RespErrorWithDate(ctx, "密码错误")
		return
	}

	ur, err := service.GetUserResource(username, s[1])
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

	ur, err := service.GetUserResource(username, filename)
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
		tool.RespErrorWithDate(ctx, "没有权限下载该文件")
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

	ctx.Writer.Header().Set("Content-Disposition", "attachment;filename="+filename)
	ctx.Writer.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	ctx.Writer.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	file.Seek(0, 0)

	for {
		var n int
		tmp := make([]byte, 10) // 通过控制切片大小控制下载速度
		n, err = file.Read(tmp)
		if err == io.EOF {
			return
		}
		ctx.Writer.Write(tmp[:n])
	}
}

func shareFile(ctx *gin.Context) {
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)
	filename := ctx.Param("filename")
	permission := ctx.Request.Header.Get("permission")

	ur, err := service.GetUserResource(username, filename)
	if err != nil {
		if err == redis.Nil {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	var str string

	if permission != "" {
		ur.Permission = permission
	}

	switch ur.Permission {
	case service.Public:
		tool.RespSuccessfulWithDate(ctx, basePath+username+"/"+filename)
	case service.Private:
		tool.RespErrorWithDate(ctx, "分享失败，您以将该文件设置为仅自己可见")
	case service.Permission:
		// 使人们只能通过分享连接下载的想法是将url进行base64编码
		str = base64.URLEncoding.EncodeToString([]byte(filename + "-" + ur.ResourceName))
		tool.RespSuccessfulWithDate(ctx, basePath+"download_conn/"+str)
	}
}

func uploadFile(ctx *gin.Context) {
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	attribute := ctx.PostForm("attribute")
	if attribute == "" {
		attribute = service.Public
	}

	folder := ctx.PostForm("folder")
	if folder == "" {
		folder = "main folder"
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		log.Println("get file failed,err:", err)
		tool.RespErrorWithDate(ctx, "上传失败")
		return
	}
	if file == nil {
		tool.RespErrorWithDate(ctx, "文件为空")
		return
	}

	res, filename, err := service.DealWithFile(file)
	if !res {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}
	if err != nil {
		tool.RespSuccessfulWithDate(ctx, err)
		return
	}

	if !service.IsRepeatFile(filename) {
		err = ctx.SaveUploadedFile(file, filename)
		if err != nil {
			log.Println("上传文件失败,err:", err)
			tool.RespErrorWithDate(ctx, "服务器错误，上传文件失败")
			return
		}
	}

	var storage = model.UserResources{
		Folder:       folder,
		Filename:     file.Filename,
		ResourceName: filename,
		Permission:   attribute,
		CreateAt:     time.Now().String(),
	}

	_, err = service.StorageFile(username, storage)
	if err != nil {
		log.Println("storage file failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
}

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
