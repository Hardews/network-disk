package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"mime/multipart"
	"network-disk/dao"
	"network-disk/model"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	Public     = "0"
	Private    = "1"
	Permission = "2"
)

var (
	ErrOfFileTooBig = errors.New("文件太大")
)

func StorageFile(username string, file *multipart.FileHeader, filename string) (bool, error) {
	var ur = model.UserResources{
		Filename:     file.Filename,
		ResourceName: filename,
		Permission:   Public,
		CreateAt:     time.Now().String(),
	}
	return dao.ResourcesFile(username, ur)
}

func GetAllUserResource(username string) ([]model.UserResources, error) {
	urMap, err := dao.GetUserAllResource(username)
	if err != nil {
		return nil, err
	}
	var urs []model.UserResources
	for filename, ur := range urMap {
		s := strings.Split(ur, ":")
		ur := model.UserResources{
			Filename:     filename,
			ResourceName: s[0],
			Permission:   s[1],
			CreateAt:     s[2],
		}
		urs = append(urs, ur)
	}
	return urs, nil
}

func GetUserResource(username, filename string) (ur model.UserResources, err error) {
	urStr, err := dao.GetUserResource(username, filename)
	if err != nil {
		return
	}

	s := strings.Split(urStr, ":")
	ur = model.UserResources{
		Filename:     filename,
		ResourceName: s[0],
		Permission:   s[1],
		CreateAt:     s[2],
	}

	return ur, nil
}

func DealWithFile(file *multipart.FileHeader) (res bool, filename string, err error) {
	res = false
	if file.Size > 1024*1024*1024*5 {
		res = true
		err = ErrOfFileTooBig
		return
	}

	fileSuffix := path.Ext(file.Filename)

	data := make([]byte, file.Size)
	dealFile, err := file.Open()
	if err != nil {
		return
	}

	_, err = dealFile.Read(data)
	if err != nil {
		return
	}

	//用md5生成唯一的文件指纹，返回保存到本地
	res = true
	filename = "./uploadFile/" + mD5(data) + fileSuffix
	return
}

func IsRepeatFile(filename string) bool {
	return filepath.IsAbs(filename)
}

func mD5(data []byte) string {
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}
