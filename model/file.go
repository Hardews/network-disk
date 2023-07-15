/**
 * @Author: Hardews
 * @Date: 2023/7/14 11:27
 * @Description:
**/

package model

import "gorm.io/gorm"

type Resource struct {
	gorm.Model
	ResourceName string // 资源名称
	ResourceNum  int    // 多少用户拥有该资源
}

type UserResources struct {
	gorm.Model
	FolderId     uint   // 外键，目标表 Folder
	ResourceId   uint   // 外键，目标表 Resource
	Filename     string // 展示的文件名称
	Permission   string // 权限
	DownloadAddr string // 下载地址
}

type Folder struct {
	gorm.Model
	Username     string // 谁创建的这个文件夹
	FolderName   string // 这个文件夹名称
	ParentFolder int    // 父文件夹 id
}
