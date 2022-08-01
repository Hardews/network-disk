package api

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	engine := gin.Default()

	engine.Use(Cors)

	engine.POST("/login", login)
	engine.POST("/register", register)

	upload := engine.Group("")
	{
		upload.Use(JwtToken)
		upload.POST("/upload", uploadFile) // 上传文件
		upload.DELETE("/upload", delFile)  // 删除文件
	}

	download := engine.Group("")
	{
		download.Use(CheckUrl)
		download.GET("/download/:filename", downloadFileByConn)        // 链接下载
		download.GET("/:username/:filename", downloadPublicFile)       // 下载公开的文件
		download.POST("/encryption/:filename", downloadEncryptionFile) // 下载加密的文件
	}

	user := engine.Group("/user")
	{
		user.Use(JwtToken)
		user.GET("/download/:filename", downloadUserFile)        // 下载用户自己的文件
		user.GET("/resource/all", getUserAllFile)                // 获取该用户的所有文件信息
		user.GET("/resource", getUserFileByCategory)             // 根据文件夹路径获取
		user.PUT("/resource", updateFileAttribute)               // 修改文件名或存储路径
		user.GET("/share/normal/:filename", shareFile)           // 正常分享文件
		user.GET("/share/QrCode/:filename", qrCode)              // 二维码分享
		user.GET("/share/encryption/:filename", encryptionShare) // 加密分享
	}

	admin := engine.Group("/admin")
	{
		admin.Use(JwtToken)
		admin.Use(AdminToken)
		admin.POST("/register", adminRegister)
		admin.GET("/resource/all", adminGetUserAllFile) // 获取用户保存的文件
		admin.PUT("/resource", adminChangeUserFile)     // 修改违禁文件
	}

	engine.Run()
}
