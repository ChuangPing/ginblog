package model

import (
	"errors"
	"fmt"
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
)

//文章表
type Article struct {
	Category Category `gorm:"foreignkey:Cid"`
	gorm.Model
	Title   string `gorm:"type:varchar(100);not null" json:"title" comment:"文章标题"`
	Cid     int    `gorm:"type:int;not null" json:"cid" comment:"文章分类id"`
	Desc    string `gorm:"type:varchar(200);comment" json:"desc" comment:"文章描述"`
	Content string `gorm:"type:longtext" json:"content" comment:"文章内容"`
	Img     string `gorm:"type:varchar(100)" json:"img" comment:"文章封面"`
}

//	CreateArticle 添加文章
func CreateArticle(article *Article) int {
	err := db.Create(&article).Error
	if err != nil {
		fmt.Println(err)
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}

// GetCateArt 查询分类下的所有文章
func GetCateArt(id int, pageSize, pageNum int) ([]Article, int, int64) {
	//	查询分类下的所有文章前先查询分类是否存在
	err := db.Where("id=?", id).First(&Category{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//	分类不存在无法修改
		return nil, errmsg.ERROR_CATEGORY_NOT_EXIST, 0
	}

	var categoryArticles []Article
	var offset int
	var total int64
	if pageNum == -1 && pageSize == -1 {
		offset = -1
	} else {
		offset = (pageNum - 1) * pageSize
	}
	err = db.Preload("Category").Where("cid=?", id).Limit(pageSize).Offset(offset).Find(&categoryArticles).Count(&total).Error
	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	return categoryArticles, errmsg.SUCCSE, total
}

//	GetArticle 查询单个文章信息
func GetArticleInfo(id int) (Article, int) {
	var article Article
	//	查询文章信息时，同时查询出关联的分类信息
	err := db.Preload("Category").Where("id=?", id).First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return article, errmsg.ERROR_ARTICLE_NOT_EXIST
	}
	return article, errmsg.SUCCSE
}

//	GetArticles 查询文章列表
func GetArticles(pageSize, pageNum int) ([]Article, int, int64) {
	var articles []Article
	var offset int
	var total int64
	//	减法保护：当用户不传page参数时，pageSize和pageNum 均为-1，此时offset的值为-1
	if pageNum == -1 && pageSize == -1 {
		offset = -1
	} else {
		offset = (pageNum - 1) * pageSize
	}
	err := db.Preload("Category").Limit(pageSize).Offset(offset).Find(&articles).Count(&total).Error
	if err != nil {
		return nil, errmsg.ERROR, 0
	}
	return articles, errmsg.SUCCSE, total
}

//	EditArticle 编辑文章
func EditArticle(id int, article *Article) int {
	//	修改文章前先查询文章是否存在
	code := CheckArticle(id)
	if code == errmsg.ERROR {
		return errmsg.ERROR_ARTICLE_NOT_EXIST
	}

	var maps = make(map[string]interface{})
	maps["title"] = article.Title
	maps["cid"] = article.Cid
	maps["desc"] = article.Desc
	maps["content"] = article.Content
	maps["img"] = article.Img
	err := db.Model(&Article{}).Where("Id=?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//	DeleteArticle 删除文章
func DeleteArticle(id int) int {
	var article Article
	//	删除文章前先查询文章是否存在
	code := CheckArticle(id)
	if code == errmsg.ERROR {
		return errmsg.ERROR_ARTICLE_NOT_EXIST
	}

	err := db.Where("id=?", id).Delete(&article).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//	CheckArticle 根据id查询文章
func CheckArticle(id int) int {
	var article Article
	err := db.Where("id=?", id).First(&article).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
