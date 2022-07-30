package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"network-disk/model"
	"network-disk/service"
	"network-disk/tool"
)

func delFile(ctx *gin.Context) {

}

func updateFileAttribute(ctx *gin.Context) {

}

func shareFile(ctx *gin.Context) {

}

func uploadFile(ctx *gin.Context) {
	var user model.User
	iUsername, _ := ctx.Get("username")
	user.Username = iUsername.(string)

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
