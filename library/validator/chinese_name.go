package validator

import (
	"github.com/dlclark/regexp2"
	"github.com/geekeryy/api-hub/library/localization"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func ChineseNameValidator(v *validator.Validate, trans ut.Translator) error {
	if err := v.RegisterValidation("chinese_name", chineseName); err != nil {
		return err
	}
	if err := v.RegisterTranslation("chinese_name", trans, func(ut ut.Translator) error {
		message := localization.Localize(ut.Locale(), &i18n.LocalizeConfig{
			MessageID:    "CHINESE_NAME_ERR_MSG",
			TemplateData: nil,
			PluralCount:  nil,
		})
		return ut.Add("chinese_name", message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		fieldName := localization.Localize(ut.Locale(), &i18n.LocalizeConfig{
			MessageID:    fe.Field(),
			TemplateData: nil,
			PluralCount:  nil,
		})
		t, _ := ut.T(fe.Tag(), fieldName)
		return t
	}); err != nil {
		return err
	}
	return nil
}

// 中文名称
var chineseNameRX = regexp2.MustCompile(`^([\u4e00-\u9fa5A-Za-z()0-9_-]){1,64}$`, regexp2.None)

func chineseName(fl validator.FieldLevel) bool {
	if v, ok := fl.Field().Interface().(string); ok {
		match, _ := chineseNameRX.MatchString(v)
		return match
	}
	return false
}
