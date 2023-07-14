/**
 * @Author: Hardews
 * @Date: 2023/7/14 12:03
 * @Description:
**/

package model

import (
	"gorm.io/gorm"
	"time"
)

type Url struct {
	gorm.Model
	Overdue time.Time
	Url     string `gorm:"not null;unique;type:varchar(200)"`
}

type Code struct {
	gorm.Model
	Url  string
	Code string
}
