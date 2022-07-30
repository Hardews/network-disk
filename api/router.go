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

	user := engine.Group("/user")
	{
		user.Use(middleware.JwtToken)
		user.GET("/resource", getUserAllFile)       // 获取该用户的所有文件信息
		user.PUT("/resource", updateFileAttribute)  // 修改文件名或存储路径
		user.GET("/share/:filename", shareFile)     // 分享文件
		user.GET("/share/QrCode/:filename", qrCode) // 二维码分享
	}

	upload := engine.Group("")
	{
		upload.Use(middleware.JwtToken)
		upload.POST("/upload", uploadFile)                         // 上传文件
		upload.DELETE("/upload", delFile)                          // 删除文件
		engine.GET("/download_conn/:filename", downloadFileByConn) // 只能通过链接下载
		engine.GET("/:username/:filename", downloadPublicFile)     // 公开的文件
	}

	engine.Run()
}
