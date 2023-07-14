package api

/*
func adminGetUserAllFile(ctx *gin.Context) {
	username, res := ctx.GetQuery("username")
	if !res {
		tool.RespErrorWithDate(ctx, "用户名为空")
		return
	}

	urs, err := service.GetAllUserResource(username)
	if err != nil {
		log.Println("admin:get user resource failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	tool.RespSuccessfulWithDate(ctx, urs)
}

// adminChangeUserFile 修改用户的文件内容(把违规的文件改为特定的文件)
func adminChangeUserFile(ctx *gin.Context) {
	filename, res := ctx.GetQuery("filename")
	if !res {
		tool.RespSuccessfulWithDate(ctx, "filename is null")
		return
	}

	folder, _ := ctx.GetQuery("folder")
	folderId, err := strconv.Atoi(folder)
	if err != nil {
		tool.RespSuccessfulWithDate(ctx, "folder id 格式错误")
		return
	}

	ur, err := service.GetUserResource(filename, folderId)
	if err != nil {
		log.Println("admin:get file info failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}

	// 删除用户原来的文件
	err = os.Remove("./uploadFile/" + service.GetResourceName(int(ur.ResourceId)))
	if err != nil {
		log.Println("admin:remove the file failed,err:", err)
		tool.RespInternetError(ctx)
		return
	}
}

*/
