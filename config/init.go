/**
 * @Author: Hardews
 * @Date: 2023/7/11 21:06
 * @Description:
**/

package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"

	encoder "github.com/zwgblue/yaml-encoder"
)

var ReloadConfig Config

type Config struct {
	BaseSetting   `yaml:"base-setting" comment:"一些基础字段"`
	MysqlDatabase `yaml:"mysql" comment:"Mysql 数据库配置"`
}

type BaseSetting struct {
	Username   string `comment:"登陆时的用户名"`
	Password   string `comment:"登陆时的密码"`
	IsSetAdmin bool   `yaml:"is-admin" comment:"为 true 时登陆账号拥有管理员权限"`
	BaseUrl    string `yaml:"url" comment:"项目运行的 url，如 127.0.0.1:8080"`
	Host       string `yaml:"host" comment:"运行时的端口号，默认为 8080"`
	LogOutAddr string `yaml:"log-out-addr" comment:"日志输出地址"`
}

// MysqlDatabase Mysql 数据库设置
type MysqlDatabase struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DatabaseLink string `yaml:"database-link"`
	DatabaseHost string `yaml:"database-host"`
	DatabaseName string `yaml:"database-name"`
}

func GenerateConfigFile() {
	file, err := os.Open("./config/config.yaml")
	defer file.Close()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("配置文件不存在...\n已自动生成，请填写后再次启动")
			err = nil
			file, err = os.Create("./config/config.yaml")
			if err != nil {
				log.Println("未知错误， err:", err)
				return
			}

			var config = Config{
				BaseSetting: BaseSetting{
					Host:       "8080",
					LogOutAddr: "./disk.log",
					IsSetAdmin: false,
				},
				MysqlDatabase: MysqlDatabase{
					Username:     "root",
					DatabaseHost: "3306",
					DatabaseName: "disk",
				},
			}
			newEncoder := encoder.NewEncoder(config, encoder.WithComments(encoder.CommentsOnHead))
			out, err := newEncoder.Encode()
			if err != nil {
				log.Println("未知错误， err:", err)
				return
			}

			file.Write(out)
			file.Close()
			os.Exit(0)
		}

		log.Println("未知错误， err:", err)
		os.Exit(1)
	}

	configByte, err := io.ReadAll(file)
	if err != nil {
		log.Println("读取配置文件错误， err:", err)
		os.Exit(1)
	}

	yaml.Unmarshal(configByte, &ReloadConfig)
}
