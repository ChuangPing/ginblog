package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

//读取配置文件信息

var (
	AppMode  string
	HttpPort string
	JwtKey   string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	AccessKey   string
	SecretKey   string
	Bucket      string
	QiniuServer string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取出错，请检查二年间路径", err)
	}
	LoadServer(file)
	LoadData(file)
	LoadQiniu(file)
}

//读取网站配置信息
func LoadServer(file *ini.File) {
	//根据ini配置文件读取数据 -- 分区，和key来读取数据，后面指定如果读取不到信息，就给变量赋默认值
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("Server").Key("JwtKey").MustString("87665J%$899")
}

//	读取数据库配置信息
func LoadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("mysql")
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassword = file.Section("database").Key("DbPassword").MustString("root")
	DbName = file.Section("database").Key("DbName").MustString("ginblog")
}

//	读取配置七牛云存储配置信息
func LoadQiniu(file *ini.File) {
	AccessKey = file.Section("qiniu").Key("AccessKey").String()
	SecretKey = file.Section("qiniu").Key("SecretKey").String()
	Bucket = file.Section("qiniu").Key("Bucket").String()
	QiniuServer = file.Section("qiniu").Key("QiniuServer").String()
}
