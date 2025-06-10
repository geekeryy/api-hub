// Package localization @Description  TODO
// @Author  	 jiangyang
// @Created  	 2024/7/18 下午5:36
package localization

import (
	"embed"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/text/language"
)

//go:embed i18n/error.*.toml
var LocaleFS embed.FS

var DefaultLanguageTag = language.Japanese

var bundle = i18n.NewBundle(DefaultLanguageTag)

func init() {
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	if _, err := bundle.LoadMessageFileFS(LocaleFS, "i18n/error.zh.toml"); err != nil {
		logx.Error(err)
	}
	if _, err := bundle.LoadMessageFileFS(LocaleFS, "i18n/error.en.toml"); err != nil {
		logx.Error(err)
	}
	if _, err := bundle.LoadMessageFileFS(LocaleFS, "i18n/error.ja.toml"); err != nil {
		logx.Error(err)
	}
	if _, err := bundle.LoadMessageFileFS(LocaleFS, "i18n/error.ko.toml"); err != nil {
		logx.Error(err)
	}
	if _, err := bundle.LoadMessageFileFS(LocaleFS, "i18n/error.zh-Hant.toml"); err != nil {
		logx.Error(err)
	}
}

func Localize(lang string, localizeConfig *i18n.LocalizeConfig) string {
	res, err := i18n.NewLocalizer(bundle, lang).Localize(localizeConfig)
	if err != nil {
		return localizeConfig.MessageID
	}
	return res
}
