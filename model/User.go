package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"ginblog/utils/errmsg"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null" json:"username" comment:"用户名" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(20);not null" json:"password" comment:"用户密码" validate:"required,min=6,max=20" label:"用户密码"`
	//	由于存在默认值，当role角色字段为0，验证器也会认为没有填而报角色字段为必填值，因此使用 1， 2 避免 0
	Role int `gorm:"type:int;default:2" json:"role" comment:"用户角色，1：管理员，2：普通用户" validate:"required" label:"角色"` // 用户角色
}

//	CheckUser 查询用户是否存在
func CheckUser(name string) (code int) {
	var user User
	err := db.Where("username=?", name).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errmsg.SUCCSE // 500
	}
	// 用户名已在数据库注册
	return errmsg.ERROR_USERNAME_USED // 1001
}

//	CreateUser 新增用户
func CreateUser(user *User) int {
	//添加用户前，先将用户密码进行加密
	user.Password = ScryptPwd(user.Password)
	err := db.Create(&user).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}

//	GetUsers 查询用户列表  --一般都会做分页处理
func GetUsers(pageSize int, pageNum int) ([]User, int, int64) {
	var users []User
	var offset int
	// 总数
	var total int64
	//	减法保护：当用户不传page参数时，pageSize和pageNum 均为-1，此时offset的值为-1
	if pageNum == -1 && pageSize == -1 {
		offset = -1
	} else {
		offset = (pageNum - 1) * pageSize
	}
	err := db.Limit(pageSize).Offset(offset).Find(&users).Count(&total).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errmsg.ERROR, 0
	}
	return users, errmsg.SUCCSE, total
}

// GetUser 根据用户id查询用户信息
func GetUser(id int) (User, int) {
	var user User
	err := db.Where("id=?", id).First(&user).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCSE
}

//	EditUser 编辑用户
func EditUser(id int, user *User) int {
	//	通过map的方式指定修改的字段
	var editMap = make(map[string]interface{})
	editMap["username"] = user.Username
	editMap["role"] = user.Role
	err := db.Model(&User{}).Where("id=?", id).Updates(editMap).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//	DeleteUser 删除用户
func DeleteUser(id int) int {
	var user User
	err := db.Where("id=?", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//	ScryptPwd 密码加密
func ScryptPwd(password string) string {
	const KeyLen = 10
	var salt = make([]byte, 8)
	salt = []byte{11, 22, 33, 44, 55, 66, 77, 88}

	hashPwd, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	//	将加密后的hash字符串经base64编码进行返回
	fpwd := base64.StdEncoding.EncodeToString(hashPwd)
	return fpwd
}

// CheckLogin 登录验证
func CheckLogin(username string, password string) int {
	var user User

	// 查询用户名是否存在
	err := db.Where("username=?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errmsg.ERROR_USER_NOT_EXIST
	}

	// 查询到的用户比对密码
	if ScryptPwd(password) != user.Password {
		return errmsg.ERROR_PASSWORD_WRONG
	}

	// 判断用户是否有权限登录后代管理系统
	if user.Role != 1 {
		return errmsg.ERROR_USER_NOT_RIGHT
	}

	return errmsg.SUCCSE
}
