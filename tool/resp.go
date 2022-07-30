package tool

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespInternetError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"info": "服务器错误",
	})
}

func RespSuccessful(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"info": "成功",
	})
}

func RespErrorWithDate(ctx *gin.Context, Data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"date": Data,
	})
}

func RespSuccessfulWithDate(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
