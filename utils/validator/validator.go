package validator

import (
	"fmt"
	"ginblog/utils/errmsg"
	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

func Validate(data interface{}) (string, int) {
	//	使用gin自带的验证器
	validate := validator.New()
	//	初始化错误信息翻译引擎  -- 使用中文的翻译引擎
	uni := ut.New(zh_Hans_CN.New())
	//	使用翻译引擎
	trans, _ := uni.GetTranslator("zh_Hans_CN")

	err := zh.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println("err", err)
	}
	//	处理label --提示错误字段使用label
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		return label
	})
	//	使用验证器验证结构的字段
	err = validate.Struct(data)
	if err != nil {
		//	验证不通过 -- 利用断言循环错误
		for _, v := range err.(validator.ValidationErrors) {
			//	错误信息进行翻译
			return v.Translate(trans), errmsg.ERROR
		}
	}
	return "", errmsg.SUCCSE
}
