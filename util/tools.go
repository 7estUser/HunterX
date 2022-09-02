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
	return strings.Split(query, "\"")[1] + "-" + RandomString(5)
}
