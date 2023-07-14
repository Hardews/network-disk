package model

import (
	"gorm.io/gorm"
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
