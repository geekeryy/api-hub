package localization

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type Language struct {
	ZH     string `json:"zh,omitempty"`
	ZHHant string `json:"zh-Hant,omitempty"`
	EN     string `json:"en,omitempty"`
	JA     string `json:"ja,omitempty"`
	KO     string `json:"ko,omitempty"`
}

const (
	ZH     = "zh"
	EN     = "en"
	JA     = "ja"
	KO     = "ko"
	ZHHant = "zh-Hant"
)
const ACCEPT_LANGUAGE = "accept-language"

func Default() string {
	return JA
}


func Lang(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		if value := ctx.Value(ACCEPT_LANGUAGE); value != nil {
			return value.(string)
		}
		return Default()
	}
	lang, ok := md[ACCEPT_LANGUAGE]
	if !ok || len(lang) == 0 || len(lang[0]) == 0 {
		return Default()
	}
	return lang[0]
}

func GetTargetLanguages(sourceLanguage string) []string {
	slice := []string{"ja", "en", "ko", "zh", "zh-Hant"}
	var result []string
	for _, s := range slice {
		if s != sourceLanguage {
			result = append(result, s)
		}
	}
	return result
}

func GetValueByLang(lang *Language, key string) string {
	result := ""
	switch key {
	case ZH:
		result = lang.ZH
	case ZHHant:
		result = lang.ZHHant
	case EN:
		result = lang.EN
	case JA:
		return lang.JA
	case KO:
		result = lang.KO
	}
	if len(result) == 0 {
		return lang.JA
	}
	return result
}

func GetValueFromJson(jsonStr, lang string) string {
	l := Unmarshal([]byte(jsonStr))
	switch lang {
	case ZH:
		return l.ZH
	case ZHHant:
		return l.ZHHant
	case EN:
		return l.EN
	case JA:
		return l.JA
	case KO:
		return l.KO
	}
	return ""
}

func ContainsEmptyValue(jsonStr string) bool {
	if jsonStr == "" {
		return false
	}
	l := Unmarshal([]byte(jsonStr))
	return anyEmpty(l.JA, l.ZH, l.EN, l.KO, l.ZHHant)
}

func anyEmpty(items ...string) bool {
	for _, item := range items {
		if item == "" {
			return true
		}
	}
	return false
}

func jsonEmpty(json string) bool {
	return json == "" || json == "{}" || json == "[]"
}

func Unmarshal(data []byte) *Language {
	if len(data) == 0 {
		return nil
	}
	if jsonEmpty(string(data)) {
		return nil
	}
	var lang Language
	err := jsoniter.Unmarshal(data, &lang)
	if err != nil {
		logx.Error(err.Error())
		return &Language{}
	}
	return &lang
}
