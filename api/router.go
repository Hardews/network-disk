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
		upload.GET("/upload", getFile)
		upload.POST("/upload", uploadFile)
		upload.PUT("/upload", updateFile)
		upload.DELETE("/upload", delFile)
	}

	engine.Run()
}
