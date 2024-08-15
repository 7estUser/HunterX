package util

import (
	"HunterX/obj"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func SelectAccountType(apiUrl string, userName string, apiKey string, proxyMe string) (string, error) {
	searchJsonData, err := SearchApi(apiUrl, userName, apiKey, "x", 1, 1, time.Now().AddDate(-1, 0, 0).Format("2006-01-02"), time.Now().Format("2006-01-02"), proxyMe)
	return searchJsonData.Data.Account_type, err
}

// hunter查询接口：分页查询
func SearchApi(apiUrl string, userName string, apiKey string, search string, page int, pageSize int, startTime string, endTime string, proxyMe string) (obj.SearchObj, error) {
	// 查询结果对象
	var searchJsonData obj.SearchObj
	// 设置代理
	var client = checkProxy(proxyMe)
	// 发送 GET 请求
	_, err := client.R().SetResult(&searchJsonData).Get(apiUrl + "/openApi/search?" +
		"username=" + userName +
		"&api-key=" + apiKey +
		"&search=" + base64UrlEncode(search) +
		"&page=" + strconv.Itoa(page) +
		"&page_size=" + strconv.Itoa(pageSize) +
		"&start_time=" + startTime +
		"&end_time=" + endTime)
	log.Printf("查询日志： " + apiUrl + "/openApi/search?" +
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
	var client = checkProxy(proxyMe)
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
		var client = checkProxy(proxyMe)
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

func DownloadFile(url string, filePath string, proxyMe string) error {
	var client = checkProxy(proxyMe)
	// 使用 Resty 客户端发送 GET 请求
	resp, err := client.R().Get(url)
	if err != nil {
		return err
	}

	// 检查响应状态码
	if resp.IsError() {
		return fmt.Errorf("failed to download file: %s", resp.Status())
	}

	// 将响应体转换为 io.Reader
	body := bytes.NewReader(resp.Body())

	// 创建文件
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 将响应体内容写入文件
	_, err = io.Copy(outFile, body)
	if err != nil {
		return err
	}

	return nil
}
