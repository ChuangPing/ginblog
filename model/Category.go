package model

import (
	"errors"
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
)

//文章分类表
type Category struct {
	ID   uint   `gorm:"primary_key;auto_increment" json:"id" comment:"文章分类id"`
	Name string `gorm:"type:varchar(20);not null" json:"name" comment:"分类名称"`
}

//	CheckCategory 检查分类是否存在
func checkCategory(name string) int {
	var category Category
	err := db.Where("name=?", name).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//	CreateCategory 创建分类
func CreateCategory(category Category) int {
	code := checkCategory(category.Name)
	if code == errmsg.SUCCSE {
		return errmsg.ERROR_CATEGORY_EXIST
	}

	err := db.Create(&category).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}

// GetCategory 查询单个分类信息
func GetCategory(id int) (Category, int) {
	var category Category
	err := db.Where("id=?", id).First(&category).Error
	if err != nil {
		return category, errmsg.ERROR
	}
	return category, errmsg.SUCCSE
}

// GetCate 获取分类列表
func GetCate(pageNum, pageSize int) ([]Category, int, int64) {
	var category []Category
	var offset int
	var total int64
	//	减法保护：当用户不传page参数时，pageSize和pageNum 均为-1，此时offset的值为-1
	if pageNum == -1 && pageSize == -1 {
		offset = -1
	} else {
		offset = (pageNum - 1) * pageSize
	}
	err := db.Limit(pageSize).Offset(offset).Find(&category).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errmsg.ERROR, 0
	}
	return category, errmsg.SUCCSE, total
}

//	EditCategory 编辑分类
func EditCategory(id int, category Category) int {
	//	修改分类前先查询分类是否存在
	err := db.Where("id=?", id).First(&Category{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//	分类不存在无法修改
		return errmsg.ERROR_CATEGORY_NOT_EXIST
	}

	var maps = make(map[string]interface{})
	maps["name"] = category.Name
	err = db.Model(Category{}).Where("id=?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// DeleteCategory 删除分类
func DeleteCategory(id int) int {
	var category Category

	err := db.Where("id=?", id).Delete(&category).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
