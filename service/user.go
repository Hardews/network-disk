package service

import (
	"errors"
	"fmt"
	"log"
	"network-disk/config"
	"network-disk/middleware"

	"gorm.io/gorm"

	"network-disk/dao"
	"network-disk/model"
)

var (
	ErrOfInternet      = errors.New("internet error")
	ErrOfNoAccount     = errors.New("账号不存在")
	ErrOfRepeatAccount = errors.New("账号已存在")
	ErrOfWrongPassword = errors.New("密码错误")
)

func GetUsernameByFolderId(folderId uint) string {
	return dao.GetUsernameByFolderId(folderId)
}

func InitUser() {
	c := config.ReloadConfig
	if c.BaseSetting.Username != "" && c.BaseSetting.Password != "" {
		register(model.User{
			Username: c.BaseSetting.Username,
			Password: c.BaseSetting.Password,
		})
		if c.BaseSetting.IsSetAdmin {
			writeAdmin(c.BaseSetting.Username)
		}
	}
}

func Login(user model.User) (res bool, token string, err error) {
	res = true
	err, flag := CheckPassword(user.Username, user.Password)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrOfNoAccount
			fmt.Println(err)
			return
		}
		err = errors.New("check password failed,err:" + err.Error())
		res = false
		return
	}

	if !flag {
		err = ErrOfWrongPassword
		return
	} else {
		err, flag = IsAdminUser(user.Username)
		if err != nil {
			log.Println("service:check admin user failed,err:", err)
		}
		if flag {
			token, flag = middleware.SetToken(user.Username, "管理员")
			if !flag {
				err = ErrOfInternet
				return
			}
		} else {
			token, flag = middleware.SetToken(user.Username, "用户")
			if !flag {
				err = ErrOfInternet
				return
			}
		}
		return
	}
}

func register(user model.User) (res bool, err error) {
	res = false
	err, flag := CheckUsername(user)
	if err != nil {
		err = errors.New("check username failed,err:" + err.Error())
		return
	}
	if !flag {
		err = ErrOfRepeatAccount
		res = true
		return
	}

	err, user.Password = Encryption(user.Password)
	if err != nil {
		return
	}

	err = WriteIn(user)
	if err != nil {
		return
	}

	res = true
	return
}

func CheckPassword(username, password string) (error, bool) {
	err, check := dao.CheckPassword(username)
	if err != nil {
		return err, false
	}
	err, res := Interpretation(check.Password, password)
	if err != nil {
		return err, false
	}
	return err, res
}

func CheckUsername(user model.User) (error, bool) {
	err := dao.CheckUsername(user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			return err, true
		}
		return err, false
	}
	return err, false
}

func WriteIn(user model.User) error {
	err := dao.WriteIn(user)
	if err != nil {
		return err
	}
	return err
}
