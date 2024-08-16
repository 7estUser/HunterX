package util

import (
	"crypto/tls"
	"encoding/base64"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/proxy"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// base64 url编码
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

func checkProxy(proxyMe string) *resty.Client {
	// 创建 REST 客户端
	var client *resty.Client
	var auth *proxy.Auth
	if proxyMe == "" {
		client = resty.New()
		// 关闭TLS验证
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		return client
	} else {
		// 代理设置
		var ipPort string
		var parts []string
		// 判断是否有用户名密码认证,设置认证信息
		if strings.Contains(proxyMe, "@") {
			// 带认证
			parts = strings.Split(proxyMe, "@")
			ipPort = parts[1]
			var newStr string
			// test:test123@127.0.0.1:7788
			if strings.HasPrefix(proxyMe, "socks5") {
				newStr = strings.ReplaceAll(parts[0], "socks5://", "")
			} else {
				newStr = strings.ReplaceAll(parts[0], "http://", "")
			}

			creds := strings.Split(newStr, ":")
			auth = &proxy.Auth{
				User:     creds[0],
				Password: creds[1],
			}
		} else {
			if strings.HasPrefix(proxyMe, "socks5") {
				ipPort = strings.ReplaceAll(proxyMe, "socks5://", "")
			} else {
				ipPort = strings.ReplaceAll(proxyMe, "http://", "")
			}
			auth = nil
		}
		// 解析代理设置
		if strings.HasPrefix(proxyMe, "socks5") {
			var dialer proxy.Dialer
			// SOCKS5代理
			//dialer, err := proxy.SOCKS5("tcp", proxyMe[len("socks5://"):], auth, proxy.Direct)
			dialer, err := proxy.SOCKS5("tcp", ipPort, auth, proxy.Direct)
			if err != nil {
				// 错误处理
				log.Printf("无法创建 SOCKS5 拨号器 : %v", err)
				dialer = nil
				os.Exit(1)
			}
			if dialer != nil {
				//log.Printf("使用 SOCKS5 代理: %v", proxyMe)
				httpClient := &http.Client{
					Transport: &http.Transport{
						Dial:            dialer.Dial,
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				}
				client = resty.NewWithClient(httpClient)
			}
		} else if strings.HasPrefix(proxyMe, "http") {
			// 解析代理地址
			proxyURL, err := url.Parse(proxyMe)
			if err != nil {
				log.Printf("解析代理地址失败 : %v", err)
			} else {
				// log.Printf("使用 HTTP 代理: %v", proxyMe)
				// 创建 HTTP Transport，并设置代理
				httpClient := &http.Client{
					Transport: &http.Transport{
						Proxy:           http.ProxyURL(proxyURL),
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				}
				client = resty.NewWithClient(httpClient)
			}
		} else {
			log.Printf("解析代理地址失败 : %v", proxyMe)
			os.Exit(1)
		}
		return client
	}
}
