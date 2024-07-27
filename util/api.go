package util

import (
	"HunterX/obj"
	"errors"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func SelectAccountType(apiUrl string, userName string, apiKey string, proxyMe string) (string, error) {
	searchJsonData, err := SearchApi(apiUrl, userName, apiKey, "x", 1, 1, time.Now().AddDate(-1, 0, 0).Format("2006-01-02"), time.Now().Format("2006-01-02"), proxyMe)
	return searchJsonData.Data.Account_type, err
}

// hunter查询接口：分页查询
func SearchApi(apiUrl string, userName string, apiKey string, search string, page int, pageSize int, startTime string, endTime string, proxyMe string) (obj.SearchObj, error) {
	// 查询结果对象
	var searchJsonData obj.SearchObj
	// 创建 REST 客户端
	var client *resty.Client
	var auth *proxy.Auth
	if proxyMe == "" {
		client = resty.New()
	} else {
		// 代理设置
		var dialer proxy.Dialer
		// 判断是否有用户名密码认证
		if strings.Contains(proxyMe, "@") {
			// 带认证
			parts := strings.Split(proxyMe, "@")
			newStr := strings.ReplaceAll(parts[0], "socks5://", "")
			creds := strings.Split(newStr, ":")
			auth = &proxy.Auth{
				User:     creds[0],
				Password: creds[1],
			}
		} else {
			auth = nil
		}
		if strings.HasPrefix(proxyMe, "socks5") {
			// SOCKS5代理
			ipPort := strings.Split(proxyMe, "@")
			//dialer, err := proxy.SOCKS5("tcp", proxyMe[len("socks5://"):], auth, proxy.Direct)
			dialer, err := proxy.SOCKS5("tcp", ipPort[1], auth, proxy.Direct)
			if err != nil {
				// 错误处理
				log.Printf("无法创建 SOCKS5 拨号器 : %v", err)
				dialer = nil
				os.Exit(1)
			}
			if dialer != nil {
				log.Printf("使用 SOCKS5 代理: %v", proxyMe)
				httpClient := &http.Client{
					Transport: &http.Transport{
						Dial: dialer.Dial,
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
				log.Printf("使用 HTTP 代理: %v", proxyMe)
				// 创建 HTTP Transport，并设置代理
				httpClient := &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyURL),
					},
				}
				client = resty.NewWithClient(httpClient)
			}
		} else {
			log.Printf("解析代理地址失败 : %v", proxyMe)
			os.Exit(1)
		}
		_ = dialer
	}

	// 发送 GET 请求
	_, err := client.R().SetResult(&searchJsonData).Get(apiUrl + "/openApi/search?" +
		"username=" + userName +
		"&api-key=" + apiKey +
		"&search=" + base64UrlEncode(search) +
		"&page=" + strconv.Itoa(page) +
		"&page_size=" + strconv.Itoa(pageSize) +
		"&start_time=" + startTime +
		"&end_time=" + endTime)
	if err != nil {
		return searchJsonData, err
	}

	if searchJsonData.Message != "success" {
		log.Println("searchApi调用错误：" + searchJsonData.Message)
	}

	return searchJsonData, nil
}

// hunter查询接口：所有结果
func SearchAllApi(apiUrl string, userName string, apiKey string, search string, startTime string, endTime string, proxyMe string) (obj.SearchObj, error) {
	var searchJsonData obj.SearchObj
	// 创建 REST 客户端
	// 创建 REST 客户端
	var client *resty.Client
	var auth *proxy.Auth
	if proxyMe == "" {
		client = resty.New()
	} else {
		// 代理设置
		var dialer proxy.Dialer
		// 判断是否有用户名密码认证
		if strings.Contains(proxyMe, "@") {
			// 带认证
			parts := strings.Split(proxyMe, "@")
			creds := strings.Split(parts[0], ":")
			auth = &proxy.Auth{
				User:     creds[0],
				Password: creds[1],
			}
		} else {
			auth = nil
		}
		if strings.HasPrefix(proxyMe, "socks5") {
			// SOCKS5代理
			ipPort := strings.Split(proxyMe, "@")
			dialer, err := proxy.SOCKS5("tcp", ipPort[1], auth, proxy.Direct)
			//dialer, err := proxy.SOCKS5("tcp", "192.168.2.189:1080", nil, proxy.Direct)
			if err != nil {
				// 错误处理
				log.Printf("无法创建 SOCKS5 拨号器 : %v", err)
				dialer = nil
				os.Exit(1)
			}
			if dialer != nil {
				log.Printf("使用 SOCKS5 代理: %v", proxyMe)
				httpClient := &http.Client{
					Transport: &http.Transport{
						Dial: dialer.Dial,
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
				log.Printf("使用 HTTP 代理: %v", proxyMe)
				// 创建 HTTP Transport，并设置代理
				httpClient := &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyURL),
					},
				}
				client = resty.NewWithClient(httpClient)
			}
		} else {
			log.Printf("解析代理地址失败 : %v", proxyMe)
			os.Exit(1)
		}
		_ = dialer
	}
	_, err := client.R().SetResult(&searchJsonData).Post(apiUrl + "/openApi/search/batch?" +
		"username=" + userName +
		"&api-key=" + apiKey +
		"&search=" + base64UrlEncode(search) +
		"&start_time=" + startTime +
		"&end_time=" + endTime)
	if searchJsonData.Message != "success" {
		log.Println("searchApi调用错误提示：" + searchJsonData.Message)
		return searchJsonData, err
	}
	log.Println("等待系统查询结果生成，请勿关闭......")
	c := 1
check:
	for true {
		// 创建 REST 客户端
		var client *resty.Client
		if proxyMe == "" {
			client = resty.New()
		} else {
			// 解析代理字符串
			parts := strings.Split(proxyMe, "@")

			// 根据代理类型创建对应的 Dialer
			var dialer proxy.Dialer
			if strings.HasPrefix(parts[len(parts)-1], "http") {
				// HTTP 代理
				proxyAddr := parts[len(parts)-1]
				auth := &proxy.Auth{}
				if len(parts) > 1 {
					creds := strings.Split(parts[0], ":")
					auth.User = creds[0]
					auth.Password = creds[1]
				}
				dialer, _ = proxy.SOCKS5("tcp", proxyAddr, auth, proxy.Direct)
			} else {
				// SOCKS5 代理
				proxyAddr := parts[len(parts)-1]
				var auth *proxy.Auth
				if len(parts) > 1 {
					creds := strings.Split(parts[0], ":")
					auth = &proxy.Auth{
						User:     creds[0],
						Password: creds[1],
					}
				}
				dialer, _ = proxy.SOCKS5("tcp", proxyAddr, auth, proxy.Direct)
			}

			// 使用代理创建 HTTP 客户端
			httpClient := &http.Client{
				Transport: &http.Transport{
					Dial: dialer.Dial,
				},
			}
			client = resty.NewWithClient(httpClient)
		}
		_, err = client.R().SetResult(&searchJsonData).Get(apiUrl + "/openApi/search/batch/" +
			strconv.Itoa(searchJsonData.Data.Task_id) +
			"?username=" + userName +
			"&api-key=" + apiKey)
		if searchJsonData.Message != "success" {
			log.Println("批量查询任务进度查询接口调用错误提示：" + searchJsonData.Message)
			return searchJsonData, err
		}
		if err != nil {
			return searchJsonData, err
		}
		if searchJsonData.Data.Progress == "100%" {
			break check
		}
		if c == 10 {
			return searchJsonData, errors.New("批量查询任务进度查询超时")
		}
		time.Sleep(time.Second * 2)
		c++
	}
	return searchJsonData, err
}

func DownloadFile(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil || resp.Body == nil {
		return err
	}
	defer resp.Body.Close()
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
