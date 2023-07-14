/**
 * @Author: Hardews
 * @Date: 2023/7/14 16:06
 * @Description:
**/

package service

import (
	"errors"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"network-disk/dao"
)

func CheckCode(url, code string) bool {
	// 也是先从 redis 拿，拿不到再从 mysql 拿
	res, err := dao.GetUrl(url)
	if errors.Is(err, redis.Nil) || errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return res == code
}

// IsOverdue 链接是否过期
func IsOverdue(url string) (bool, error) {
	// 从 redis 拿，拿不到从 mysql 拿
	_, err := dao.GetUrl(url)
	if err == nil {
		return true, nil
	} else if err == redis.Nil || err == gorm.ErrRecordNotFound {
		return false, nil
	} else {
		return false, err
	}
}

func SetExpirationTime(url string, et int, code ...string) error {
	if len(code) != 0 {
		return dao.SetExpirationTime(url, et, code[0])
	}
	return dao.SetExpirationTime(url, et)
}
