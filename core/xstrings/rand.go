package xstrings

import "math/rand"

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

func GenerateRandomNumber(n int) string {
	numbers := "0123456789"

	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = numbers[rand.Int63()%int64(len(numbers))]
	}

	return string(result)
}
