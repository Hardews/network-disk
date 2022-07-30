package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"io"
	"log"
	"net/http"
	"network-disk/model"
	"network-disk/service"
	"network-disk/tool"
	"os"
	"strconv"
	"strings"
)

func delFile(ctx *gin.Context) {

}

func updateFileAttribute(ctx *gin.Context) {

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

	switch ur.Permission {
	case service.Public:
		tool.RespSuccessfulWithDate(ctx, "http://127.0.0.1:8080/"+username+"/"+filename)
	case service.Private:
		tool.RespErrorWithDate(ctx, "分享失败，您以将该文件设置为仅自己可见")
	case service.Permission:
		// 使人们只能通过分享连接下载的想法是将url进行base64编码
		str = base64.URLEncoding.EncodeToString([]byte(filename + "-" + ur.ResourceName))
		tool.RespSuccessfulWithDate(ctx, "http://127.0.0.1:8080/download_conn/"+str)
	}
}

func uploadFile(ctx *gin.Context) {
	var user model.User
	iUsername, _ := ctx.Get("username")
	user.Username = iUsername.(string)

	attribute := ctx.PostForm("attribute")
	if attribute == "" {
		attribute = service.Public
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

	res, err = service.StorageFile(user.Username, file, filename)
	if err != nil {
		log.Println("storage file failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}
	if res {
		tool.RespErrorWithDate(ctx, "上传失败,文件已存在")
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
