/**
 * @Author: Hardews
 * @Date: 2023/7/14 15:16
 * @Description:
**/

package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
	"net/http"

	"log"
	"strconv"

	"network-disk/service"
	"network-disk/tool"
)

// CheckUrl 检查连接是否过期或是否存在
func CheckUrl(ctx *gin.Context) {
	res, err := service.IsOverdue(ctx.Request.RequestURI)
	if err != nil {
		log.Println("upload:check due failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}
	if !res {
		tool.RespErrorWithDate(ctx, "链接无效或已过期")
		ctx.Abort()
		return
	}
}

// shareFile 分享文件
func shareFile(ctx *gin.Context) {
	// 获取文件名
	filename, exists := ctx.GetQuery("filename")
	if !exists {
		tool.RespSuccessfulWithDate(ctx, "filename 为空")
		return
	}

	// 获取 folder id
	folder, _ := ctx.GetQuery("folder")
	folderId, err := strconv.Atoi(folder)
	if err != nil {
		tool.RespErrorWithDate(ctx, "param error")
		return
	}

	// 获取过期时间
	et, _ := ctx.GetQuery("time")
	expirationTime, err := strconv.Atoi(et)
	if err != nil {
		log.Println("upload:translate et failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	// 获取该用户想分享的资源的信息
	ur, err := service.GetUserResource(filename, folderId)
	if err != nil {
		if errors.Is(err, redis.Nil) || errors.Is(err, gorm.ErrRecordNotFound) {
			tool.RespErrorWithDate(ctx, "没有该文件")
			return
		}
		log.Println(err)
		tool.RespInternetError(ctx)
		return
	}

	var url string
	switch ur.Permission {
	case service.Public:
		url = basePath + "download?encryption=false&folder=" + strconv.Itoa(folderId) + "&filename=" + filename
	case service.Private:
		tool.RespErrorWithDate(ctx, "分享失败，您以将该文件设置为仅自己可见")
	}
	err = service.SetExpirationTime(url, expirationTime)
	if err != nil {
		log.Println("upload:set et failed,err:", err)
	}
	tool.RespSuccessfulWithDate(ctx, url)
}

func encryptionShare(ctx *gin.Context) {
	pwd, res := ctx.GetQuery("code")
	var tmp string
	if !res {
		// 如果没指定就自己生成 code
		pwd = service.RandomStr(4)
		tmp = pwd
	}
	pwd = service.MD5([]byte(pwd))

	filename, res := ctx.GetQuery("filename")
	if !res {
		tool.RespSuccessfulWithDate(ctx, "filename is null")
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

	folder, _ := ctx.GetQuery("folder")
	folderId, err := strconv.Atoi(folder)
	if err != nil {
		tool.RespSuccessfulWithDate(ctx, "folder id 格式错误")
		return
	}

	url := basePath + "download?encryption=true&folder=" + strconv.Itoa(folderId) + "&filename=" + filename

	// 写入时间
	err = service.SetExpirationTime(url, eTime, tmp)
	if err != nil {
		log.Println("upload:set et failed,err:", err)
		tool.RespErrorWithDate(ctx, "分享失败,请重试")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": tmp,
		"path": url,
	})
}

// qrCode 二维码分享，用了github.com/skip2/go-qrcode包
func qrCode(ctx *gin.Context) {
	filename, res := ctx.GetQuery("filename")
	if !res {
		tool.RespSuccessfulWithDate(ctx, "filename is null")
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

	folder, _ := ctx.GetQuery("folder")
	folderId, err := strconv.Atoi(folder)
	if err != nil {
		tool.RespSuccessfulWithDate(ctx, "folder id 格式错误")
		return
	}

	ur, err := service.GetUserResource(filename, folderId)
	if err != nil {
		if errors.Is(err, redis.Nil) || errors.Is(err, gorm.ErrRecordNotFound) {
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
		url = basePath + "download?encryption=false&folder=" + strconv.Itoa(folderId) + "&filename=" + filename
	case service.Private:
		tool.RespErrorWithDate(ctx, "分享失败，您以将该文件设置为仅自己可见")
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
