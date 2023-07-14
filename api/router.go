package api

import (
	"github.com/gin-gonic/gin"
	"network-disk/config"
	"network-disk/middleware"
)

func InitRouter() {
	engine := gin.Default()

	engine.Use(middleware.Cors)

	engine.POST("/token", login) // 获取凭证

	upload := engine.Group("")
	{
		upload.Use(middleware.JwtToken)
		upload.POST("/upload", uploadFile) // 上传文件
		upload.DELETE("/upload", delFile)  // 删除文件
	}

	download := engine.Group("")
	{
		download.Use(CheckUrl)
		// 下载 baseUrl/download?encryption=bool&folder=?&filename=?
		download.GET("/download", downloadFileByConn)
	}

	user := engine.Group("/posts")
	{
		user.Use(middleware.JwtToken)
		user.GET("/download/:folderId", downloadUserFile) // 下载用户自己的文件

		folder := user.Group("/folder")
		{
			folder.GET("", getFolderInfo) // 获取该用户所有文件夹信息
			folder.POST("", addFolder)    // 添加文件夹
		}

		resource := user.Group("/resource")
		{
			resource.GET("/all", getUserAllFile)    // 获取该用户的所有文件信息
			resource.GET("", getUserFileByCategory) // 根据文件夹路径获取
			// resource.PUT("", updateFileAttribute)   // 修改文件名或存储路径 TODO 待开发
		}

		share := user.Group("/share")
		{
			share.GET("/normal", shareFile)           // 正常分享文件
			share.GET("/QrCode", qrCode)              // 二维码分享
			share.GET("/encryption", encryptionShare) // 加密分享
		}
	}

	/*
		弃用，（只有自己用）
		admin := engine.Group("/admin")
		{
			admin.Use(middleware.JwtToken)
			admin.GET("/resource/all", adminGetUserAllFile) // 获取用户保存的文件
			admin.PUT("/resource", adminChangeUserFile)     // 修改违禁文件
		}

	*/

	var host string
	if config.ReloadConfig.Host != "" {
		host = ":" + config.ReloadConfig.Host
	}
	engine.Run(host)
}
