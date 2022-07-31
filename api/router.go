package api

import (
	"github.com/gin-gonic/gin"

	"network-disk/middleware"
)

func InitRouter() {
	engine := gin.Default()

	engine.Use(middleware.Cors)

	engine.POST("/login", login)
	engine.POST("/register", register)

	upload := engine.Group("")
	{
		upload.Use(middleware.JwtToken)
		upload.POST("/upload", uploadFile)                           // 上传文件
		upload.DELETE("/upload", delFile)                            // 删除文件
		upload.GET("/download/:filename", downloadFileByConn)        // 只能通过链接下载
		upload.GET("/:username/:filename", downloadPublicFile)       // 公开的文件
		upload.POST("/encryption/:filename", downloadEncryptionFile) // 加密的文件
	}

	user := engine.Group("/user")
	{
		user.Use(middleware.JwtToken)
		user.GET("/resource/all", getUserAllFile)                // 获取该用户的所有文件信息
		user.GET("/resource", getUserFileByCategory)             // 根据文件夹路径获取
		user.PUT("/resource", updateFileAttribute)               // 修改文件名或存储路径
		user.GET("/share/normal/:filename", shareFile)           // 正常分享文件
		user.GET("/share/QrCode/:filename", qrCode)              // 二维码分享
		user.GET("/share/encryption/:filename", encryptionShare) // 加密分享
	}

	engine.Run()
}
