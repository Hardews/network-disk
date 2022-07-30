package main

import (
	"network-disk/api"
	"network-disk/dao"
)

func main() {
	dao.InitDB()
	api.InitRouter()
}
