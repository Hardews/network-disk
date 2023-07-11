package api

import (
	"github.com/gin-gonic/gin"

	"log"
	"os"

	"network-disk/model"
	"network-disk/service"
	"network-disk/tool"
)

func adminGetUserAllFile(ctx *gin.Context) {
	username, res := ctx.GetQuery("username")
	if !res {
		tool.RespErrorWithDate(ctx, "用户名为空")
		return
	}

	urs, err := service.GetAllUserResource(username)
	if err != nil {
		log.Println("admin:get user resource failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	tool.RespSuccessfulWithDate(ctx, urs)
}

// adminChangeUserFile 修改用户的文件内容(把违规的文件改为特定的文件)
func adminChangeUserFile(ctx *gin.Context) {
	changeFile, err := ctx.FormFile("file")
	if err != nil {
		log.Println("admin:change file failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}
	// 获取文件名
	filename := ctx.PostForm("filename")

	// 获取用户名
	username := ctx.PostForm("username")

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

	ur, err := service.GetUserResource(username, filename, Path, folder)
	if err != nil {
		log.Println("admin:get file info failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	// 删除用户原来的文件
	err = os.Remove(ur.ResourceName)
	if err != nil {
		log.Println("admin:remove the file failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	// 按源路径保存要更换的文件
	err = ctx.SaveUploadedFile(changeFile, ur.ResourceName)
	if err != nil {
		log.Println("admin:save the file failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}
}

func adminRegister(ctx *gin.Context) {
	var user model.User
	user.Username, _ = ctx.GetPostForm("username")
	user.Password, _ = ctx.GetPostForm("password")
	if user.Username == "" {
		tool.RespErrorWithDate(ctx, "用户名为空")
		return
	}
	if user.Password == "" {
		tool.RespErrorWithDate(ctx, "密码为空")
		return
	}

	res, err := service.Register(user)
	if !res {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}
	if err != nil {
		tool.RespErrorWithDate(ctx, err.Error())
		return
	}

	err = service.WriteAdmin(user.Username)
	if err != nil {
		log.Println("admin:write admin failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
}
