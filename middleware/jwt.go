package middleware

import (
	"fmt"
	"ginblog/utils"
	"ginblog/utils/errmsg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var JwtKey = []byte(utils.JwtKey)

//	MyClaims 生成jwt的自定义字段
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// SetToken	生成token
func SetToken(username string) (string, int) {
	//	设置过期时间  10小时
	expireTime := time.Now().Add(10 * time.Hour)
	//	初始化自定义结构体  -- jwt根据这个结构生成对应的jwt
	SetClaim := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			// token颁发签名
			Issuer: "chuanPing",
		},
	}
	//	使用指定的签名方法创建签名对象
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaim)
	//	使用指定的Secret签名并获得完整的编码后的字符串token  -- 进行签名
	token, err := reqClaim.SignedString(JwtKey)
	if err != nil {
		return "", errmsg.ERROR
	}
	return token, errmsg.SUCCSE
}

//	CheckToken 验证token
func CheckToken(tokenStr string) (*MyClaims, int) {
	//	解析token
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, errmsg.ERROR
	}

	//	验证令牌是否生效
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, errmsg.SUCCSE
	} else {
		return nil, errmsg.ERROR
	}
}

//	JwtToken jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 	从请求头中获取token
		tokenHeader := c.Request.Header.Get("Authorization")
		//fmt.Println("tokenHeader", tokenHeader)
		if tokenHeader == "" {
			code := errmsg.ERROR_TOKEN_EXIST
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			// 	阻止 中间件向下运行，不能使用return而要使用abort ，将指针置为最大
			c.Abort()
			//	阻止下面代码继续运行  -- 这两个都不能少
			return
		}
		cToken := strings.SplitN(tokenHeader, " ", 2)
		//fmt.Println("cToken", cToken)
		// "Bearer" 是固定的写法
		if len(cToken) != 2 && cToken[0] != "Bearer" {
			code := errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
		}
		claims, tCode := CheckToken(cToken[1])
		if tCode == errmsg.ERROR {
			code := errmsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			// 阻止代码继续运行
			c.Abort()
			// 阻止代码向下运行
			return
		}

		// 验证通过判断是否在有效期
		if time.Now().Unix() > claims.ExpiresAt {
			code := errmsg.ERROR_TOKEN_RUNTIME
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		//	token有效且未过期验证通过
		c.Set("username", claims.Username)
		c.Next()
	}
}
