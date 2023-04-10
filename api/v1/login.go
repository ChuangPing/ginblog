package v1

import (
	"fmt"
	"ginblog/middleware"
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println("登录接口，接受用户信息失败", err)
		return
	}
	var token string
	var code int
	code = model.CheckLogin(user.Username, user.Password)

	if code == errmsg.SUCCSE {
		// 为登录成功的用户颁发token
		token, code = middleware.SetToken(user.Username)
	}

	// 错误状态均在code，统一返回
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"token":   token,
	})
}
