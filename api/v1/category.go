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

//	AddCategory 添加分类
func AddCategory(c *gin.Context) {
	var category model.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		fmt.Println("绑定category出错", err)
		return
	}

	code = model.CreateCategory(category)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    category,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetCategoryInfo 查询单个分类
func GetCategoryInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println("GetCategoryInfo:读出参数id出错")
		return
	}

	//	根据id查询分类信息
	category, code := model.GetCategory(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    category,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetCate 查询分类列表
func GetCate(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	if pageNum == 0 {
		// pageNum = -1 根据gorm可以查询全部
		pageNum = -1
	}
	if pageSize == 0 {
		pageSize = -1
	}

	users, code, total := model.GetCate(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    users,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

//	EditCategory 编辑分类
func EditCategory(c *gin.Context) {
	//	获取编辑用户id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal("编辑分类接口接受用户参数出错", err)
	}
	//	获取用户编辑更新内容
	var category model.Category
	err = c.ShouldBindJSON(&category)
	if err != nil {
		log.Fatal("编辑分类接口接受用户参数出错", err)
	}
	code = model.EditCategory(id, category)
	if code == errmsg.ERROR_USERNAME_USED {
		//	阻止gin向下运行
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

//	删除分类
func DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal("删除分类接口接受用户id出错", err)
	}
	code = model.DeleteCategory(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
