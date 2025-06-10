// Package validator @Description
// @Author  	 jiangyang
// @Created  	 2024/6/5 下午2:23
package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/ko"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

type SliceValidationError []error

// Error concatenates all error elements in SliceValidationError into a single string separated by \n.
func (err SliceValidationError) Error() string {
	n := len(err)
	switch n {
	case 0:
		return ""
	default:
		var b strings.Builder
		if err[0] != nil {
			fmt.Fprintf(&b, "[%d]: %s", 0, err[0].Error())
		}
		// TODO 改为warp不用\n
		if n > 1 {
			for i := 1; i < n; i++ {
				if err[i] != nil {
					b.WriteString("\n")
					fmt.Fprintf(&b, "[%d]: %s", i, err[i].Error())
				}
			}
		}
		return b.String()
	}
}

var _ binding.StructValidator = (*defaultValidator)(nil)

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		if value.Elem().Kind() != reflect.Struct {
			return v.ValidateStruct(value.Elem().Interface())
		}
		return v.validateStruct(obj)
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *defaultValidator) validateStruct(obj any) error {
	v.lazyinit()
	err := v.validate.Struct(obj)

	if verrs, okV := err.(validator.ValidationErrors); _uni != nil && okV && len(verrs) > 0 {
		errStr := make([]string, 0, len(verrs))
		if t, found := _uni.GetTranslator("zh"); found {
			for _, verr := range verrs {
				errStr = append(errStr, verr.Translate(t))
			}
		}
		// TODO 改为warp不用\n
		return errors.New(strings.Join(errStr, "\n"))
	}

	return err

}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://pkg.go.dev/github.com/go-playground/validator/v10
func (v *defaultValidator) Engine() any {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		if _validator == nil {
			v.validate = validator.New()
		} else {
			v.validate = _validator
		}
		v.validate.SetTagName("binding")
	})
}

type ValidatorFn func(v *validator.Validate, trans ut.Translator) error

var _validator *validator.Validate
var _once sync.Once
var _uni *ut.UniversalTranslator

func Register(fns []ValidatorFn, langs []string) {
	_once.Do(func() {
		v := validator.New()
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
			for _, fn := range fns {
				fn(v, trans)
			}
		}

		_validator = v

		// 替换gin默认验证器
		binding.Validator = &defaultValidator{}

	})

}
