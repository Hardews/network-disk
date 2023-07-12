package service

import (
	"gorm.io/gorm"
	"network-disk/dao"
)

func IsAdminUser(username string) (err error, res bool) {
	res = false
	err = dao.CheckAdminUser(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			return
		}
		return
	}
	return nil, true
}

func writeAdmin(username string) error {
	return dao.WriteAdmin(username)
}
