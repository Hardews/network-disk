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
		tool.RespErrorWithDate(ctx, err)
	}

	tool.RespSuccessfulWithDate(ctx, token)
}
func register(ctx *gin.Context) {
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
		tool.RespInternetError(ctx)
		log.Println(err)
		return
	}
	if err != nil {
		tool.RespErrorWithDate(ctx, err)
		return
	}

	tool.RespSuccessful(ctx)
}
