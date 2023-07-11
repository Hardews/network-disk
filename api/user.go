package api

import (
	"log"

	"network-disk/model"
	"network-disk/service"
	"network-disk/tool"

	"github.com/gin-gonic/gin"
)

func login(ctx *gin.Context) {
	var user model.User
	var res bool
	user.Username, res = ctx.GetPostForm("username")
	if !res {
		tool.RespErrorWithDate(ctx, "输入的账号为空")
		return
	}

	user.Password, res = ctx.GetPostForm("password")
	if !res {
		tool.RespErrorWithDate(ctx, "输入的密码为空")
		return
	}

	res, token, err := service.Login(user)
	if !res {
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}
	if err != nil {
		tool.RespErrorWithDate(ctx, err.Error())
		return
	}

	tool.RespSuccessfulWithDate(ctx, token)
}
