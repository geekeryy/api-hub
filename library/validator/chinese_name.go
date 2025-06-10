package validator

import (
	"github.com/dlclark/regexp2"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func ChineseNameValidator(v *validator.Validate, trans ut.Translator) error {
	if err := v.RegisterValidation("chinese_name", chineseName); err != nil {
		panic(err)
	}
	if err := v.RegisterTranslation("chinese_name", trans, func(ut ut.Translator) error {
		// TODO 翻译 text
		return ut.Add("chinese_name", "{0}不符合中文名称规范，可包含中文、大小写字母、数字、_、-、()，且长度介于1到64个字符之间", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		// TODO 翻译 field
		t, _ := ut.T("chinese_name", fe.Field())
		return t
	}); err != nil {
		panic(err)
	}
	return nil
}

// 中文名称
var chineseNameRX = regexp2.MustCompile(`^([\u4e00-\u9fa5A-Za-z()0-9_-]){1,64}$`, regexp2.None)

func chineseName(fl validator.FieldLevel) bool {
	if v, ok := fl.Field().Interface().(string); ok {
		match,_:= chineseNameRX.MatchString(v)
		return match
	}
	return false
}
