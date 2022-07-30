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
		user.GET("/resource", getUserAllFile)      // 获取该用户的所有文件信息
		user.PUT("/resource", updateFileAttribute) // 修改文件名或存储路径
		user.GET("/share/:filename", shareFile)
	}

	upload := engine.Group("/upload")
	{
		upload.Use(middleware.JwtToken)
		upload.POST("/", uploadFile) // 上传文件
		upload.DELETE("/", delFile)  // 删除文件
	}

	engine.Run()
}
