package consts

type contextKey string

const AcceptLanguage contextKey = "accept-language"

func (c contextKey) String() string {
	return string(c)
}
