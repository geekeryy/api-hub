package xstrings

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHash 加密密码 PasswordHash(PlainText)
func PasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// PasswordMatch 匹配密码 PasswordMatch(PasswordHash(PlainText),PlainText)
func PasswordMatch(pwdFromDB, pwdFromInput string) bool {
	return bcrypt.CompareHashAndPassword([]byte(pwdFromDB), []byte(pwdFromInput)) == nil
}
