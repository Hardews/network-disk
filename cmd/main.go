package main

import (
	"network-disk/api"
	"network-disk/config"
	"network-disk/dao"
	"network-disk/service"
	"network-disk/tool"
)

func main() {
	config.GenerateConfigFile()
	tool.InitLog()
	dao.InitDB()
	service.InitUser()
	api.InitRouter()
}
