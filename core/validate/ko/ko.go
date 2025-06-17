package ko

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/geekeryy/api-hub/library/localization"
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// RegisterDefaultTranslations registers a set of default translations
// for all built in tag's in validator; you may add your own as desired.
func RegisterDefaultTranslations(v *validator.Validate, trans ut.Translator) (err error) {

	translations := []struct {
		tag             string
		translation     string
		override        bool
		customRegisFunc validator.RegisterTranslationsFunc
		customTransFunc validator.TranslationFunc
	}{
		{
			tag:         "required",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "required_if",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "required_unless",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "required_with",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "required_with_all",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "required_without",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "required_without_all",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "excluded_if",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "excluded_unless",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "excluded_with",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "excluded_with_all",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "excluded_without",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "excluded_without_all",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag:         "isdefault",
			translation: "{0}은(는) 필수 필드입니다",
			override:    false,
		},
		{
			tag: "len",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("len-string", "{0}은(는) {1}자 이상이어야 합니다", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("len-string-character", "{0}字符", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("len-string-character", "{0}자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("len-number", "{0}은(는) {1}이어야 합니다", false); err != nil {
					return
				}

				if err = ut.Add("len-items", "{0}은(는) {1}개 이상이어야 합니다", false); err != nil {
					return
				}
				//if err = ut.AddCardinal("len-items-item", "{0}项", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("len-items-item", "{0}개", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("len-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("len-string", translateFieldName(ut, fe), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("len-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("len-items", translateFieldName(ut, fe), c)

				default:
					t, err = ut.T("len-number", translateFieldName(ut, fe), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("警告: 翻译字段错误: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "min",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("min-string", "{0}은(는) {1}자 이상이어야 합니다", false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0}글자", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.Add("min-number", "{0}은(는) {1}이어야 합니다", false); err != nil {
					return
				}

				if err = ut.Add("min-items", "{0}은(는) {1}개 이상이어야 합니다", false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-items-item", "{0}개", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("min-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("min-string", translateFieldName(ut, fe), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("min-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("min-items", translateFieldName(ut, fe), c)

				default:
					t, err = ut.T("min-number", translateFieldName(ut, fe), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("警告: 翻译字段错误: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "max",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("max-string", "{0}은(는) {1}자 이하이어야 합니다", false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0}글자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("max-number", "{0}은(는) {1}이어야 합니다", false); err != nil {
					return
				}

				if err = ut.Add("max-items", "{0}은(는) {1}개 이하이어야 합니다", false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-items-item", "{0}개", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("max-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("max-string", translateFieldName(ut, fe), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("max-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("max-items", translateFieldName(ut, fe), c)

				default:
					t, err = ut.T("max-number", translateFieldName(ut, fe), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("警告: 翻译字段错误: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eq",
			translation: "{0}不等于{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					fmt.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ne",
			translation: "{0}不能等于{1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					fmt.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "lt",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("lt-string", "{0}은(는) {1}자 이하이어야 합니다", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0}글자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-number", "{0}은(는) {1}이어야 합니다", false); err != nil {
					return
				}

				if err = ut.Add("lt-items", "{0}은(는) {1}개 이하이어야 합니다", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0}개", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-datetime", "{0}은(는) 현재 날짜와 시간보다 이전이어야 합니다", false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lt-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-string", translateFieldName(ut, fe), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lt-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-items", translateFieldName(ut, fe), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s'不能用于struct类型.", fe.Tag())
					} else {
						t, err = ut.T("lt-datetime", translateFieldName(ut, fe))
					}

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-number", translateFieldName(ut, fe), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("警告: 翻译字段错误: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "lte",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("lte-string", "{0}은(는) {1}자 이하이어야 합니다", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0}글자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-number", "{0}은(는) {1}이어야 합니다", false); err != nil {
					return
				}

				if err = ut.Add("lte-items", "{0}은(는) {1}개 이하이어야 합니다", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0}개", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-datetime", "{0}은(는) 현재 날짜와 시간보다 이전이어야 합니다", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lte-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-string", translateFieldName(ut, fe), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lte-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-items", translateFieldName(ut, fe), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s'은(는) struct 타입에 사용할 수 없습니다", fe.Tag())
					} else {
						t, err = ut.T("lte-datetime", translateFieldName(ut, fe))
					}

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-number", translateFieldName(ut, fe), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("警告: 翻译字段错误: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gt",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("gt-string", "{0}은(는) {1}자 이상이어야 합니다", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("gt-string-character", "{0}个字符", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("gt-string-character", "{0}글자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-number", "{0}은(는) {1}이어야 합니다", false); err != nil {
					return
				}

				if err = ut.Add("gt-items", "{0}은(는) {1}개 이상이어야 합니다", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("gt-items-item", "{0}项", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("gt-items-item", "{0}개", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-datetime", "{0}은(는) 현재 날짜와 시간보다 이후이어야 합니다", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gt-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-string", translateFieldName(ut, fe), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gt-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-items", translateFieldName(ut, fe), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s'은(는) struct 타입에 사용할 수 없습니다", fe.Tag())
					} else {

						t, err = ut.T("gt-datetime", translateFieldName(ut, fe))
					}

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-number", translateFieldName(ut, fe), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("警告: 翻译字段错误: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gte",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("gte-string", "{0}은(는) {1}자 이상이어야 합니다", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("gte-string-character", "{0}个字符", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("gte-string-character", "{0}글자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-number", "{0}은(는) {1}이어야 합니다", false); err != nil {
					return
				}

				if err = ut.Add("gte-items", "{0}은(는) {1}개 이상이어야 합니다", false); err != nil {
					return
				}

				//if err = ut.AddCardinal("gte-items-item", "{0}项", locales.PluralRuleOne, false); err != nil {
				//	return
				//}

				if err = ut.AddCardinal("gte-items-item", "{0}개", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-datetime", "{0}은(는) 현재 날짜와 시간보다 이후이어야 합니다", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gte-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-string", translateFieldName(ut, fe), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gte-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-items", translateFieldName(ut, fe), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s'은(는) struct 타입에 사용할 수 없습니다", fe.Tag())
					} else {
						t, err = ut.T("gte-datetime", translateFieldName(ut, fe))
					}

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-number", translateFieldName(ut, fe), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("警告: 翻译字段错误: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqfield",
			translation: "{0}은(는) {1}이어야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqcsfield",
			translation: "{0}은(는) {1}이어야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "necsfield",
			translation: "{0}은(는) {1}이 아니어야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtcsfield",
			translation: "{0}은(는) {1}보다 크야야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtecsfield",
			translation: "{0}은(는) {1}보다 크거나 같아야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltcsfield",
			translation: "{0}은(는) {1}보다 작아야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltecsfield",
			translation: "{0}은(는) {1}보다 작거나 같아야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "nefield",
			translation: "{0}은(는) {1}이 아니어야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtfield",
			translation: "{0}은(는) {1}보다 크야야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtefield",
			translation: "{0}은(는) {1}보다 크거나 같아야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltfield",
			translation: "{0}은(는) {1}보다 작아야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltefield",
			translation: "{0}은(는) {1}보다 작거나 같아야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "alpha",
			translation: "{0}은(는) 영문자만 포함할 수 있습니다",
			override:    false,
		},
		{
			tag:         "alphanum",
			translation: "{0}은(는) 영문자와 숫자만 포함할 수 있습니다",
			override:    false,
		},
		{
			tag:         "alphanumunicode",
			translation: "{0}은(는) 영문자, 숫자 및 Unicode 문자만 포함할 수 있습니다",
			override:    false,
		},
		{
			tag:         "alphaunicode",
			translation: "{0}은(는) 영문자와 Unicode 문자만 포함할 수 있습니다",
			override:    false,
		},
		{
			tag:         "numeric",
			translation: "{0}은(는) 유효한 숫자여야 합니다",
			override:    false,
		},
		{
			tag:         "number",
			translation: "{0}은(는) 유효한 숫자여야 합니다",
			override:    false,
		},
		{
			tag:         "hexadecimal",
			translation: "{0}은(는) 유효한 16진수여야 합니다",
			override:    false,
		},
		{
			tag:         "hexcolor",
			translation: "{0}은(는) 유효한 16진수 색상이어야 합니다",
			override:    false,
		},
		{
			tag:         "rgb",
			translation: "{0}은(는) 유효한 RGB 색상이어야 합니다",
			override:    false,
		},
		{
			tag:         "rgba",
			translation: "{0}은(는) 유효한 RGBA 색상이어야 합니다",
			override:    false,
		},
		{
			tag:         "hsl",
			translation: "{0}은(는) 유효한 HSL 색상이어야 합니다",
			override:    false,
		},
		{
			tag:         "hsla",
			translation: "{0}은(는) 유효한 HSLA 색상이어야 합니다",
			override:    false,
		},
		{
			tag:         "email",
			translation: "{0}은(는) 유효한 이메일 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "url",
			translation: "{0}은(는) 유효한 URL이어야 합니다",
			override:    false,
		},
		{
			tag:         "uri",
			translation: "{0}은(는) 유효한 URI이어야 합니다",
			override:    false,
		},
		{
			tag:         "base64",
			translation: "{0}은(는) 유효한 Base64 문자열이어야 합니다",
			override:    false,
		},
		{
			tag:         "contains",
			translation: "{0}은(는) {1}을(를) 포함해야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "containsany",
			translation: "{0}은(는) 최소한 {1} 중 하나를 포함해야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "containsrune",
			translation: "{0}은(는) {1}을(를) 포함해야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludes",
			translation: "{0}은(는) {1}을(를) 포함하지 않아야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesall",
			translation: "{0}은(는) {1}을(를) 포함하지 않아야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesrune",
			translation: "{0}은(는) {1}을(를) 포함하지 않아야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "endswith",
			translation: "{0}은(는) {1}로 끝나야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "startswith",
			translation: "{0}은(는) {1}로 시작해야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "isbn",
			translation: "{0}은(는) 유효한 ISBN 번호여야 합니다",
			override:    false,
		},
		{
			tag:         "isbn10",
			translation: "{0}은(는) 유효한 ISBN-10 번호여야 합니다",
			override:    false,
		},
		{
			tag:         "isbn13",
			translation: "{0}은(는) 유효한 ISBN-13 번호여야 합니다",
			override:    false,
		},
		{
			tag:         "issn",
			translation: "{0}은(는) 유효한 ISSN 번호여야 합니다",
			override:    false,
		},
		{
			tag:         "uuid",
			translation: "{0}은(는) 유효한 UUID여야 합니다",
			override:    false,
		},
		{
			tag:         "uuid3",
			translation: "{0}은(는) 유효한 V3 UUID여야 합니다",
			override:    false,
		},
		{
			tag:         "uuid4",
			translation: "{0}은(는) 유효한 V4 UUID여야 합니다",
			override:    false,
		},
		{
			tag:         "uuid5",
			translation: "{0}은(는) 유효한 V5 UUID여야 합니다",
			override:    false,
		},
		{
			tag:         "ulid",
			translation: "{0}은(는) 유효한 ULID여야 합니다",
			override:    false,
		},
		{
			tag:         "ascii",
			translation: "{0}은(는) ASCII 문자만 포함할 수 있습니다",
			override:    false,
		},
		{
			tag:         "printascii",
			translation: "{0}은(는) 인쇄 가능한 ASCII 문자만 포함할 수 있습니다",
			override:    false,
		},
		{
			tag:         "multibyte",
			translation: "{0}은(는) 다중 바이트 문자를 포함해야 합니다",
			override:    false,
		},
		{
			tag:         "datauri",
			translation: "{0}은(는) 유효한 데이터 URI를 포함해야 합니다",
			override:    false,
		},
		{
			tag:         "latitude",
			translation: "{0}은(는) 유효한 위도 좌표를 포함해야 합니다",
			override:    false,
		},
		{
			tag:         "longitude",
			translation: "{0}은(는) 유효한 경도 좌표를 포함해야 합니다",
			override:    false,
		},
		{
			tag:         "ssn",
			translation: "{0}은(는) 유효한 소득 원천 증명서(SSN)여야 합니다",
			override:    false,
		},
		{
			tag:         "ipv4",
			translation: "{0}은(는) 유효한 IPv4 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "ipv6",
			translation: "{0}은(는) 유효한 IPv6 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "ip",
			translation: "{0}은(는) 유효한 IP 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "cidr",
			translation: "{0}은(는) 유효한 무클래스 도메인 간 라우팅(CIDR)여야 합니다",
			override:    false,
		},
		{
			tag:         "cidrv4",
			translation: "{0}은(는) 유효한 IPv4 주소를 포함하는 무클래스 도메인 간 라우팅(CIDR)여야 합니다",
			override:    false,
		},
		{
			tag:         "cidrv6",
			translation: "{0}은(는) 유효한 IPv6 주소를 포함하는 무클래스 도메인 간 라우팅(CIDR)여야 합니다",
			override:    false,
		},
		{
			tag:         "tcp_addr",
			translation: "{0}은(는) 유효한 TCP 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "tcp4_addr",
			translation: "{0}은(는) 유효한 IPv4 TCP 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "tcp6_addr",
			translation: "{0}은(는) 유효한 IPv6 TCP 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "udp_addr",
			translation: "{0}은(는) 유효한 UDP 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "udp4_addr",
			translation: "{0}은(는) 유효한 IPv4 UDP 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "udp6_addr",
			translation: "{0}은(는) 유효한 IPv6 UDP 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "ip_addr",
			translation: "{0}은(는) 유효한 IP 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "ip4_addr",
			translation: "{0}은(는) 유효한 IPv4 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "ip6_addr",
			translation: "{0}은(는) 유효한 IPv6 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "unix_addr",
			translation: "{0}은(는) 유효한 UNIX 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "mac",
			translation: "{0}은(는) 유효한 MAC 주소여야 합니다",
			override:    false,
		},
		{
			tag:         "iscolor",
			translation: "{0}은(는) 유효한 색상이어야 합니다",
			override:    false,
		},
		{
			tag:         "oneof",
			translation: "{0}은(는) {1} 중 하나여야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				s, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}
				return s
			},
		},
		{
			tag:         "json",
			translation: "{0}은(는) JSON 문자열이어야 합니다",
			override:    false,
		},
		{
			tag:         "lowercase",
			translation: "{0}은(는) 소문자여야 합니다",
			override:    false,
		},
		{
			tag:         "uppercase",
			translation: "{0}은(는) 대문자여야 합니다",
			override:    false,
		},
		{
			tag:         "datetime",
			translation: "{0}의 형식은 {1}이어야 합니다",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), translateFieldName(ut, fe), fe.Param())
				if err != nil {
					log.Printf("警告: 翻译字段错误: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "image",
			translation: "{0}은(는) 유효한 이미지여야 합니다",
			override:    false,
		},
	}

	for _, t := range translations {

		if t.customTransFunc != nil && t.customRegisFunc != nil {

			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)

		} else if t.customTransFunc != nil && t.customRegisFunc == nil {

			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)

		} else if t.customTransFunc == nil && t.customRegisFunc != nil {

			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFunc)

		} else {
			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		}

		if err != nil {
			return
		}
	}

	return
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {

	return func(ut ut.Translator) (err error) {

		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return

	}

}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	fieldName := localization.Localize(ut.Locale(), &i18n.LocalizeConfig{
		MessageID:    translateFieldName(ut, fe),
		TemplateData: nil,
		PluralCount:  nil,
	})
	t, err := ut.T(fe.Tag(), fieldName)
	if err != nil {
		log.Printf("警告: 翻译字段错误: %#v", fe)
		return fe.(error).Error()
	}

	return t
}

func translateFieldName(ut ut.Translator, fe validator.FieldError) string {
	fieldName := localization.Localize(ut.Locale(), &i18n.LocalizeConfig{
		MessageID:    fe.Field(),
		TemplateData: nil,
		PluralCount:  nil,
	})
	return fieldName
}
