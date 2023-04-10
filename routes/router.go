package routes

import (
	v1 "ginblog/api/v1"
	"ginblog/middleware"
	"ginblog/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	//	设置gin运行状态：生产或者开发环境
	gin.SetMode(utils.AppMode)
	r := gin.New()
	r.Use(gin.Recovery())
	// 使用自定义日志
	r.Use(middleware.Logger())
	// 使用跨域中间件
	r.Use(middleware.Cors())
	auth := r.Group("api/v1")
	//	需要权限的路由  -- 验证token
	auth.Use(middleware.JwtToken())
	{
		//	user模块的路由接口

		//	编辑用户路由
		auth.PUT("user/:id", v1.EditUser)
		//	删除用户路由
		auth.DELETE("user/:id", v1.DeleteUser)

		//	分类模块的路由接口
		//	添加分类路由
		auth.POST("category/add", v1.AddCategory)
		//	编辑分类路由
		auth.PUT("category/:id", v1.EditCategory)
		//	删除分类路由
		auth.DELETE("category/:id", v1.DeleteCategory)

		//	文章模块的路由接口
		//	添加文章路由
		auth.POST("article/add", v1.AddArticle)
		//	编辑文章路由
		auth.PUT("article/:id", v1.EditArticle)
		//	删除文章路由
		auth.DELETE("article/:id", v1.DeleteArticle)

		//	上串文件路由
		auth.POST("upload", v1.Upload)
	}
	// GET请求均不需要鉴权
	router := r.Group("api/v1")
	{
		//	添加用户路由  -- 不需要权限
		router.POST("user/add", v1.AddUser)
		//	查询用户列表
		router.GET("users", v1.GetUsers)
		//	获取单个用户信息
		router.GET("user/:id", v1.GetUserInfo)
		//	获取分类列表
		router.GET("category", v1.GetCate)
		// 	获取文章列表
		router.GET("article/list", v1.GetArticles)
		//	获取分类下的所有文章
		router.GET("article/:id", v1.GetCateArt)
		//	获取单个文章
		router.GET("article/info/:id", v1.GetArticleInfo)
		//	用户登录路由
		router.POST("login", v1.Login)
	}
	r.Run(utils.HttpPort)
}
