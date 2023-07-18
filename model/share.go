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
	gorm.Model `json:"base_info,omitempty"`
	Overdue    time.Time `json:"overdue,omitempty"`
	Url        string    `gorm:"not null;type:varchar(200)"`
}

type Code struct {
	gorm.Model `json:"base_info,omitempty"`
	Url        string
	Code       string
}
