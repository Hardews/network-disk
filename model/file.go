/**
 * @Author: Hardews
 * @Date: 2023/7/14 11:27
 * @Description:
**/

package model

import "gorm.io/gorm"

type Resource struct {
	gorm.Model   `json:"base_info,omitempty" `
	ResourceName string `json:"resource_name,omitempty"` // 资源名称
	ResourceNum  int    `json:"resource_num,omitempty"`  // 多少用户拥有该资源
}

type UserResources struct {
	gorm.Model   `json:"base_info,omitempty"`
	FolderId     uint   `json:"folder_id,omitempty"`     // 外键，目标表 Folder
	ResourceId   uint   `json:"resource_id,omitempty"`   // 外键，目标表 Resource
	Filename     string `json:"filename,omitempty"`      // 展示的文件名称
	Permission   string `json:"permission,omitempty"`    // 权限
	DownloadAddr string `json:"download_addr,omitempty"` // 下载地址
}

type Folder struct {
	gorm.Model   `json:"base_info,omitempty"`
	Username     string `json:"username"`      // 谁创建的这个文件夹
	FolderName   string `json:"folder_name"`   // 这个文件夹名称
	ParentFolder int    `json:"parent_folder"` // 父文件夹 id
}
