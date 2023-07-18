package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"base_info,omitempty"`
	Username   string `gorm:"not null;type:varchar(20)"`
	Password   string `gorm:"type:varchar(255)"`
}

type AdminUser struct {
	gorm.Model `json:"base_info,omitempty"`
	Username   string
}
