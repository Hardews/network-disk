package main

import (
	"network-disk/api"
	"network-disk/dao"
	"network-disk/tool"
)

func main() {
	tool.InitLog()
	dao.InitDB()
	api.InitRouter()
}
