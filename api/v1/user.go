package v1

import (
	"fmt"
	"ginblog/model"
	"ginblog/utils/errmsg"
	"ginblog/utils/validator"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//	定义全局状态码变量
var code int

//	查询用户是否存在
func UserExist(c *gin.Context) {

}

//	AddUser 添加用户
func AddUser(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		fmt.Println("绑定user出错", err)
		return
	}
	//	验证前端传递数据是否合法
	msg, code := validator.Validate(&user)
	if code != errmsg.SUCCSE {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": msg,
		})
		return
	}
	//	添加用户前先检查用户名是否注册
	code = model.CheckUser(user.Username)
	if code == errmsg.SUCCSE {
		code = model.CreateUser(&user)
	}

	//	统一处理错误，包括 CheckUser 和 CreateUser，错误码都通过code进行接受
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    user,
		"message": errmsg.GetErrMsg(code),
	})
}

//	GetUserInfo 查询单个用户
func GetUserInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal("查询单个用户获取id出错", err)
	}

	//	根据id查询用户信息
	user, code := model.GetUser(id)
	var maps = make(map[string]interface{})
	maps["username"] = user.Username
	maps["role"] = user.Role
	//	这里统一返回，因为如果有错误会包含在code中，向前端返回的data为空，前端通过判断status就行判断
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    maps,
		"message": errmsg.GetErrMsg(code),
	})
}

//	GetUsers 查询用户列表
func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	//if err != nil {
	//	fmt.Println("获取数据出错", err)
	//	return
	//}
	if pageNum == 0 {
		// pageNum = -1 根据gorm可以查询全部
		pageNum = -1
	}
	if pageSize == 0 {
		pageSize = -1
	}

	users, code, total := model.GetUsers(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    users,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

//	EditUser 编辑用户
func EditUser(c *gin.Context) {
	//	获取编辑用户id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// log.Fatal() 类似于 log.Print() 函数，后跟调用 os.Exit(1) 函数  --打印后退出程序，不用在写return
		log.Fatal("编辑用户接口接受用户参数出错", err)
	}
	//	获取用户编辑更新内容
	var user model.User
	err = c.ShouldBindJSON(&user)
	if err != nil {
		// log.Fatal() 类似于 log.Print() 函数，后跟调用 os.Exit(1) 函数  --打印后退出程序，不用在写return
		log.Fatal("编辑用户接口接受用户参数出错", err)
	}
	fmt.Println(id, user)

	//	更新前用户修改的用户是否重名
	code = model.CheckUser(user.Username)
	if code == errmsg.SUCCSE {
		code = model.EditUser(id, &user)
	}
	if code == errmsg.ERROR_USERNAME_USED {
		//	阻止gin向下运行
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})

}

//	DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	// DELETE /api/v1/user/:id   -- localhost:3000/api/v1/user/1  获取路径上这种方式传递的参数使用 c.Param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal("删除用户接口接受用户id出错", err)
	}
	code = model.DeleteUser(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
