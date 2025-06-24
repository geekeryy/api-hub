// Package validator @Description
// @Author  	 jiangyang
// @Created  	 2024/6/5 下午2:23
package validate

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/geekeryy/api-hub/core/language"
	en_translations "github.com/geekeryy/api-hub/core/validate/en"
	ja_translations "github.com/geekeryy/api-hub/core/validate/ja"
	ko_translations "github.com/geekeryy/api-hub/core/validate/ko"
	zh_translations "github.com/geekeryy/api-hub/core/validate/zh"
	zh_Hant_translations "github.com/geekeryy/api-hub/core/validate/zh_hant"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/ko"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
)

type Validate struct {
	validate *validator.Validate
}

type ValidatorFn func(v *validator.Validate, trans ut.Translator) error

var _validate *Validate
var _once sync.Once
var _uni *ut.UniversalTranslator

func New(fns []ValidatorFn, langs []string) *Validate {
	_once.Do(func() {
		v := validator.New(validator.WithRequiredStructEnabled())
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return fld.Tag.Get("comment")
		})
		_uni = ut.New(zh.New())
		for _, lang := range langs {
			var t locales.Translator
			switch lang {
			case "zh":
				t = zh.New()
			case "en":
				t = en.New()
			case "ja":
				t = ja.New()
			case "ko":
				t = ko.New()
			case "zh_Hant":
				t = zh_Hant.New()
			}
			_uni.AddTranslator(t, true)
			trans, ok := _uni.GetTranslator(t.Locale())
			if !ok {
				continue
			}
			switch lang {
			case "zh":
				zh_translations.RegisterDefaultTranslations(v, trans)
			case "en":
				en_translations.RegisterDefaultTranslations(v, trans)
			case "ja":
				ja_translations.RegisterDefaultTranslations(v, trans)
			case "ko":
				ko_translations.RegisterDefaultTranslations(v, trans)
			case "zh_Hant":
				zh_Hant_translations.RegisterDefaultTranslations(v, trans)
			}
			for _, fn := range fns {
				if err := fn(v, trans); err != nil {
					logx.Error(err)
				}
			}
		}

		_validate = &Validate{
			validate: v,
		}

	})
	return _validate

}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *Validate) ValidateStruct(ctx context.Context, obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		if value.Elem().Kind() != reflect.Struct {
			return v.ValidateStruct(ctx, value.Elem().Interface())
		}
		return v.validateStruct(ctx, obj)
	case reflect.Struct:
		return v.validateStruct(ctx, obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		// 如果验证失败，返回第一个错误
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(ctx, value.Index(i).Interface()); err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *Validate) validateStruct(ctx context.Context, obj any) error {
	err := v.validate.Struct(obj)

	// 如果验证失败，返回第一个错误
	if verrs, okV := err.(validator.ValidationErrors); _uni != nil && okV && len(verrs) > 0 {
		lang := language.Lang(ctx)
		if t, found := _uni.GetTranslator(lang); found {
			for _, verr := range verrs {
				return errors.New(verr.Translate(t))
			}
		}
	}

	return err
}
