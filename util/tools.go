package util

import (
	"encoding/base64"
	"math/rand"
	"strings"
)

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

//生成随机字符串
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune
	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//base64 url编码
func base64UrlEncode(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

func OutFileName(query string) string {
	startIndex := strings.Index(query, "=\"")
	endIndex := strings.Index(query[startIndex+1:], "\"")
	randomString := generateRandomString(4)
	if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
		return query[startIndex+1:endIndex] + "-" + randomString
	} else {
		return randomString + "-" + randomString
	}
}

// 生成指定长度的随机字符数
func generateRandomString(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
