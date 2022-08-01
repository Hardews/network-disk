package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"

	"network-disk/dao"
	"network-disk/model"
)

const (
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

func IsOverdue(url string) (bool, error) {
	_, err := dao.GetUrl(url)
	if err == nil {
		return true, nil
	} else {
		if err == redis.Nil {
			err = nil
		} else {
			return false, err
		}
	}

	err = dao.GetForeverUrl(url)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func SetExpirationTime(url string, et int) error {
	return dao.SetExpirationTime(url, et)
}

func GetUserFileByCategory(username string, category string, Path string) ([]model.UserResources, error) {
	urs, err := GetAllUserResource(username)
	if err != nil {
		return nil, err
	}

	var res []model.UserResources
	for _, ur := range urs {
		if ur.Folder == category && ur.Path == Path {
			res = append(res, ur)
		}
	}

	return res, nil
}

// UpdateFileAttribute 更新文件属性
func UpdateFileAttribute(old model.UserResources, new, username string, chose int) (res bool, err error) {
	// 先删除后更改
	_, err = dao.DelResourceFile(username, old.Filename, old.Path, old.Folder)
	if err != nil {
		return
	}
	switch chose {
	case updateName:
		old.Filename = new
		var urs []model.UserResources
		urs, err = GetAllUserResource(username)
		if err != nil {
			return
		}
		for _, ur := range urs {
			// 判断是否重命名
			if ur.Filename == new && ur.Folder == old.Folder && ur.Path == old.Path {
				return false, ErrOfSameName
			}
		}
	case updateAttr:
		old.Permission = new
	case updatePath:
		old.Path = new
	default:
		err = ErrOfNoKnow
		return
	}
	return dao.ResourcesFile(username, old)
}

// DelFile 删除文件
func DelFile(username string, filename, resource, path, folder string) (err error) {
	// 检查该用户是否有存储该文件
	_, err = dao.GetUserResource(username, filename, path, folder)
	if err != nil {
		return
	}

	n, err := dao.ResourceDecr(filename)
	if err != nil {
		return
	}
	if n <= 0 {
		// 没有人存储这个文件了，删除
		err = os.Remove(resource)
		if err != nil {
			return
		}
	}

	_, err = dao.DelResourceFile(username, filename, path, folder)
	return
}

func StorageFile(username string, ur model.UserResources) (bool, error) {
	_, err := dao.ResourceIncr(ur.ResourceName)
	if err != nil {
		return false, err
	}
	res, err := dao.ResourcesFile(username, ur)
	return res, err
}

func GetAllUserResource(username string) ([]model.UserResources, error) {
	urMap, err := dao.GetUserAllResource(username)
	if err != nil {
		return nil, err
	}

	// 遍历得到的map，处理获取到的数据
	var urs []model.UserResources
	for key, ur := range urMap {
		s1 := strings.Split(key, "&&")
		s2 := strings.Split(ur, "&&")
		ur := model.UserResources{
			Path:     s1[0],
			Filename: s1[1],
			Folder:   s1[2],

			ResourceName: s2[0],
			Permission:   s2[1],
			CreateAt:     s2[2],
			DownloadAddr: s2[3],
		}
		urs = append(urs, ur)
	}
	return urs, nil
}

// GetUserResource 根据信息来获取资源信息
func GetUserResource(username, filename, path, folder string) (ur model.UserResources, err error) {
	urStr, err := dao.GetUserResource(username, filename, path, folder)
	if err != nil {
		return
	}

	s := strings.Split(urStr, "&&")
	ur = model.UserResources{
		Filename:     filename,
		ResourceName: s[0],
		Permission:   s[1],
		CreateAt:     s[2],
		Folder:       folder,
		DownloadAddr: s[3],
		Path:         path,
	}

	return ur, nil
}

func IsRepeatFilename(username, filename, folder, path string) (res bool, err error) {
	urs, err := GetAllUserResource(username)
	if err != nil {
		return
	}

	for _, ur := range urs {
		if ur.Filename == filename && ur.Folder == folder && ur.Path == path {
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
	return nil
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

	// 判断是否存在这个文件后缀的文件夹
	bathPath := "./uploadFile/"
	_, err = os.Stat(bathPath + fileSuffix[1:])
	if err == nil {
		// 存在则存入对应文件夹
		bathPath += fileSuffix[1:] + "/"
	} else if os.IsNotExist(err) {
		// 不存在则默认路径
		err = nil
	}

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
