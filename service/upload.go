package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"network-disk/dao"
	"network-disk/model"
	"os"
	"path"
	"strconv"
)

const (
	bathPath   = "./uploadFile/"
	Public     = "0"
	Private    = "1"
	Permission = "2"
	updateName = 1
	updateAttr = 2
	updatePath = 3
)

var (
	ErrOfFileTooBig = errors.New("文件太大")
	ErrOfNoKnow     = errors.New("不懂你想干嘛")
	ErrOfSameName   = errors.New("文件名重复")
)

func StorageFile(ur model.UserResources) (bool, error) {
	// 该资源拥有者数量增加
	_, err := dao.ResourceIncr(ur.ResourceId)
	if err != nil {
		return false, err
	}
	// 存入数据库
	res, err := dao.ResourcesFile(ur)
	return res, err
}

func IsRepeatFilename(filename string, folderId int) (res bool, err error) {
	urs, err := GetAllUserResource(dao.GetUsernameByFolderId(uint(folderId)))
	if err != nil {
		return
	}

	for _, ur := range urs {
		if ur.Filename == filename && ur.ID == uint(folderId) {
			return false, nil
		}
	}
	return true, nil
}

func Storage(file *multipart.FileHeader, resourceName, breakPointPath string, breakFile *os.File) (err error) {
	_, err = os.Stat(breakPointPath)
	if !os.IsNotExist(err) {
		// 如果不存在断点文件则创建
		breakFile, err = os.Create(breakPointPath)
		defer breakFile.Close()
		if err != nil {
			err = errors.New("create break point file failed,err:" + err.Error())
			return
		}
	}

	// 读取断点位置
	breakFile, err = os.OpenFile(breakPointPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		err = errors.New("open break point file failed,err:" + err.Error())
		return
	}
	b, err := ioutil.ReadAll(breakFile)
	if err != nil {
		return
	}
	start := string(b)

	var resourceFile *os.File
	fmt.Println(resourceName)
	resourceFile, err = os.OpenFile(resourceName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		} else {
			err = errors.New("open resource file failed,err:" + err.Error())
			return
		}
	}

	upFile, err := file.Open()
	if err != nil {
		err = errors.New("open upload file failed,err:" + err.Error())
		return
	}

	// 存储
	count, _ := strconv.ParseInt(start, 10, 64)
	resourceFile.Seek(count, 0)
	upFile.Seek(count, 0)
	data := make([]byte, 1024, 1024)
	var upTotal, total, Len = 0, 0, 0

	for {
		total, err = upFile.Read(data)
		if err == io.EOF {
			// 删除文件 需要先关闭该文件
			err = upFile.Close()
			err = resourceFile.Close()
			err = breakFile.Close()
			if err != nil {
				err = errors.New("临时记录文件关闭失败" + err.Error())
				log.Println(err)
			}
			err = os.Remove(breakPointPath)
			if err != nil {
				err = errors.New("临时记录文件删除失败" + err.Error())
				log.Println(err)
			}
			break
		}
		Len, err = resourceFile.Write(data[:total])
		if err != nil {
			err = errors.New("write file failed,err:" + err.Error())
			return
		}
		upTotal += Len
		// 记录上传长度
		count += int64(Len)
		breakFile.Seek(0, 0)
		breakFile.WriteString(strconv.Itoa(int(count)))
		// 模拟断开
		//if count > 4438903 {
		//  log.Fatal("模拟上传中断")
		//}
	}
	// 到这里时上传结束了，更新 resource 表
	return dao.CreateANewResource(resourceName)
}

// DealWithFile 对文件预处理
func DealWithFile(file *multipart.FileHeader) (res bool, filename string, err error) {
	res = false
	// 判断文件是否过大
	if file.Size > 1024*1024*1024*5 {
		res = true
		err = ErrOfFileTooBig
		return
	}

	// 获取后缀
	fileSuffix := path.Ext(file.Filename)

	// 直接存在一个大文件夹
	// 读取文件内容
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
	filename = bathPath + MD5(data) + fileSuffix
	return
}

func MD5(data []byte) string {
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}
