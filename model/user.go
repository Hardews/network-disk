package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;type:varchar(20)"`
	Password string `gorm:"type:varchar(255)"`
}

type AdminUser struct {
	gorm.Model
	Username string
}
