package middleware

import (
	"github.com/gin-gonic/gin"

	"log"

	"network-disk/service"
	"network-disk/tool"
)

func CheckUrl(ctx *gin.Context) {
	res, err := service.IsOverdue(ctx.Request.URL.String())
	if err != nil {
		log.Println("upload:check due failed,err:", err)
		return
	}
	if !res {
		tool.RespErrorWithDate(ctx, "链接无效或已过期")
		return
	}
}
