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
	mainFolder         = "主文件夹"
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
		_, flag := CheckUsername(model.User{Username: c.BaseSetting.Username})
		if flag && c.BaseSetting.IsSetAdmin {
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

	// 账号密码写入
	err = dao.WriteIn(user)
	if err != nil {
		return
	}

	// 创建默认的信息，比如主文件夹
	_, err = dao.CreateFolder(model.Folder{
		Username:     user.Username,
		FolderName:   mainFolder,
		ParentFolder: -1,
	})
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
