package xstrings

import (
	"math/rand"

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

func GenerateRandomString(n int) string {
	letters := "ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz"
	numbers := "123456789"

	result := make([]byte, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			result[i] = letters[rand.Int63()%int64(len(letters))]
		} else {
			result[i] = numbers[rand.Int63()%int64(len(numbers))]
		}
	}

	return string(result)
}
