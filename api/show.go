/**
 * @Author: Hardews
 * @Date: 2023/7/18 17:00
 * @Description:
**/

package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"log"
	"network-disk/service"
	"network-disk/tool"
	"os"
	"strconv"
	"strings"
)

const (
	picture = iota
)

func isSupportFile(fileType int, filename string) bool {
	sufArr := strings.Split(filename, ".")
	if len(sufArr) < 2 {
		return false
	}
	switch fileType {
	case picture:
		if sufArr[len(sufArr)-1] == "jpg" || sufArr[len(sufArr)-1] == "png" {
			return true
		}
		return false
	default:
		return false
	}
}

func showPicture(ctx *gin.Context) {
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

	// 判断是不是图片
	if !isSupportFile(picture, filename) {
		tool.RespSuccessfulWithDate(ctx, "不支持的文件类型")
		return
	}

	ur, err := service.GetUserResource(filename, folderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, redis.Nil) {
			tool.RespSuccessfulWithDate(ctx, "没找到该文件")
			return
		}
		log.Println("get user resource error,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	rName := service.GetResourceName(int(ur.ResourceId))
	picFile, err := os.Open(rName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, redis.Nil) {
			tool.RespSuccessfulWithDate(ctx, "没找到该文件")
			return
		}
		log.Println("open user resource error,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	var res []byte
	for true {
		var temp = make([]byte, 1024)
		n, err := picFile.Read(temp)
		if err != nil {
			log.Println("read user resource error,err:", err)
			tool.RespInternetError(ctx)
			return
		}
		res = append(res, temp[:n]...)
	}
	ctx.Writer.Write(res)
}
