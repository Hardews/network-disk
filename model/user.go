package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique;type:varchar(20)"`
	Password string `gorm:"type:varchar(100)"`
}

type AdminUser struct {
	gorm.Model
	Username string `gorm:"not null;unique;type:varchar(20)"`
}

type Url struct {
	gorm.Model
	Overdue time.Time
	Url     string `gorm:"not null;unique;type:varchar(200)"`
}

type UserResources struct {
	gorm.Model
	Username     string
	Folder       string
	Path         string
	Filename     string
	ResourceName string
	Permission   string
	DownloadAddr string
	CreateAt     string
}
