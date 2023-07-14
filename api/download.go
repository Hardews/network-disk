/**
 * @Author: Hardews
 * @Date: 2023/7/14 15:13
 * @Description:文件下载相关
**/

package api

import (
	"github.com/gin-gonic/gin"
	"network-disk/service"

	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"network-disk/tool"
)

func downloadFileByConn(ctx *gin.Context) {
	encryption, exists := ctx.GetQuery("encryption")
	if !exists {
		tool.RespErrorWithDate(ctx, "encryption 为空")
		return
	}
	isEncryption, err := strconv.ParseBool(encryption)
	if err != nil {
		tool.RespSuccessfulWithDate(ctx, "加密字段传入错误")
		return
	}
	if isEncryption {
		// 加密分享
		code, exists := ctx.GetQuery("code")
		if !exists {
			tool.RespSuccessfulWithDate(ctx, "请输入提取码！！")
			return
		}
		if !service.CheckCode(ctx.Request.URL.String(), code) {
			tool.RespSuccessfulWithDate(ctx, "提取码错误！！")
			return
		}
	}

	// 获取 folder id
	folder, exists := ctx.GetQuery("folder")
	if !exists {
		tool.RespErrorWithDate(ctx, "folder 为空")
		return
	}
	folderId, err := strconv.Atoi(folder)
	if err != nil {
		tool.RespErrorWithDate(ctx, "folderId 为空")
		return
	}

	filename, exists := ctx.GetQuery("filename")
	if !exists {
		tool.RespErrorWithDate(ctx, "filename 为空")
		return
	}

	ur, err := service.GetUserResource(filename, folderId)
	if err != nil {
		tool.RespInternetError(ctx)
		return
	}

	downloadFile(ctx, ur.Filename, service.GetResourceName(int(ur.ResourceId)))
}

// downloadUserFile 用户下载自己的资源
func downloadUserFile(ctx *gin.Context) {
	// 获取用户名
	iUsername, _ := ctx.Get("username")
	username := iUsername.(string)

	folder := ctx.Param("folderId")
	folderId, err := strconv.Atoi(folder)
	if err != nil {
		tool.RespErrorWithDate(ctx, "folderId 为空")
		return
	}

	filename, exists := ctx.GetQuery("filename")
	if !exists {
		tool.RespErrorWithDate(ctx, "filename 为空")
		return
	}

	if username != service.GetUsernameByFolderId(uint(folderId)) {
		ctx.JSON(403, "Forbidden!")
		return
	}

	ur, err := service.GetUserResource(filename, folderId)
	if err != nil {
		tool.RespInternetError(ctx)
		return
	}

	downloadFile(ctx, ur.Filename, service.GetResourceName(int(ur.ResourceId)))

}

// downloadFile 文件下载
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

	// 设置响应头
	ctx.Writer.Header().Set("Content-Disposition", "attachment;filename="+filename)
	ctx.Writer.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	ctx.Writer.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	// 设置初始偏移量
	file.Seek(0, 0)

	for {
		var n int
		// 通过控制切片大小控制下载速度
		tmp := make([]byte, 1024*10) // 10 兆
		n, err = file.Read(tmp)
		if err == io.EOF {
			return
		}

		ctx.Writer.Write(tmp[:n])
		time.Sleep(1 * time.Millisecond)
	}
}
