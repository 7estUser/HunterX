package util

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

//base64 url编码
func base64UrlEncode(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

func OutFileName(query string) string {
	startIndex := strings.Index(query, "=\"")
	endIndex := strings.Index(query[startIndex+2:], "\"")
	randomString := generateRandomString(4)
	if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
		if endIndex > 10 {
			return query[startIndex+2:startIndex+2+9] + "-" + randomString + ".xlsx"
		}
		return query[startIndex+2:startIndex+2+endIndex] + "-" + randomString + ".xlsx"
	} else {
		return randomString + "-" + randomString + ".xlsx"
	}
}

// 生成指定长度的随机字符数
func generateRandomString(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
