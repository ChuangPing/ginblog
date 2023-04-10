package v1

import (
	"fmt"
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//	添加文章
func AddArticle(c *gin.Context) {
	var article model.Article
	err := c.ShouldBindJSON(&article)
	if err != nil {
		fmt.Println("绑定user出错", err)
		return
	}
	code = model.CreateArticle(&article)

	//	统一处理错误，包括 CheckUser 和 CreateUser，错误码都通过code进行接受
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    article,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetCateArt 查询分类下的所有文章
func GetCateArt(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal("查询分类下的所有文章获取分类id出错", err)
	}
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	if pageNum == 0 {
		// pageNum = -1 根据gorm可以查询全部
		pageNum = -1
	}
	if pageSize == 0 {
		pageSize = -1
	}
	categoryArticles, code, total := model.GetCateArt(id, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    categoryArticles,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

//	GetArticleInfo 查询单个文章
func GetArticleInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal("查询单个文章获取id出错", err)
	}

	//	根据id查询文章信息
	article, code := model.GetArticleInfo(id)
	//	这里统一返回，因为如果有错误会包含在code中，向前端返回的data为空，前端通过判断status就行判断
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    article,
		"message": errmsg.GetErrMsg(code),
	})
}

// SearchArticle 搜索文章标题 TODO

//	查询文章列表
func GetArticles(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	if pageNum == 0 {
		// pageNum = -1 根据gorm可以查询全部
		pageNum = -1
	}
	if pageSize == 0 {
		pageSize = -1
	}

	users, code, total := model.GetArticles(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    users,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

//	EditArticle 编辑文章
func EditArticle(c *gin.Context) {
	//	获取编辑用户id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// log.Fatal() 类似于 log.Print() 函数，后跟调用 os.Exit(1) 函数  --打印后退出程序，不用在写return
		log.Fatal("编辑文章接口接受用户参数出错", err)
	}
	//	获取用户编辑更新内容
	var article model.Article
	err = c.ShouldBindJSON(&article)
	if err != nil {
		// log.Fatal() 类似于 log.Print() 函数，后跟调用 os.Exit(1) 函数  --打印后退出程序，不用在写return
		log.Fatal("编辑文章接口接受用户参数出错", err)
	}
	code = model.EditArticle(id, &article)
	if code == errmsg.ERROR_USERNAME_USED {
		//	阻止gin向下运行
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

//	DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal("删除文章接口接受用户id出错", err)
	}
	code = model.DeleteArticle(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
