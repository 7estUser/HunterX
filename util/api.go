package util

import (
	"HunterX/obj"
	"errors"
	"github.com/go-resty/resty/v2"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func SelectAccountType(userName string, apiKey string) (string, error) {
	searchJsonData, err := SearchApi(userName, apiKey, "x", 1, 1, time.Now().AddDate(-1, 0, 0).Format("2006-01-02"), time.Now().Format("2006-01-02"))
	return searchJsonData.Data.Account_type, err
}

//hunter查询接口：分页查询
func SearchApi(userName string, apiKey string, search string, page int, pageSize int, startTime string, endTime string) (obj.SearchObj, error) {
	//查询结果对象：Obj
	var searchJsonData obj.SearchObj
	//创建client对象
	client := resty.New()
	//调用searchApi Get请求接口
	println("https://hunter.qianxin.com/openApi/search?" +
		"username=" + userName +
		"&api-key=" + apiKey +
		"&search=" + base64UrlEncode(search) +
		"&page=" + strconv.Itoa(page) +
		"&page_size=" + strconv.Itoa(pageSize) +
		"&start_time=" + startTime +
		"&end_time=" + endTime)
	_, err := client.R().SetResult(&searchJsonData).Get("https://hunter.qianxin.com/openApi/search?" +
		"username=" + userName +
		"&api-key=" + apiKey +
		"&search=" + base64UrlEncode(search) +
		"&page=" + strconv.Itoa(page) +
		"&page_size=" + strconv.Itoa(pageSize) +
		"&start_time=" + startTime +
		"&end_time=" + endTime)
	if searchJsonData.Message != "success" {
		log.Fatalln("searchApi调用错误：" + searchJsonData.Message)
	}
	return searchJsonData, err
}

//hunter查询接口：所有结果
func SearchAllApi(userName string, apiKey string, search string, startTime string, endTime string) (obj.SearchObj, error) {
	//查询结果Obj
	var searchJsonData obj.SearchObj
	client := resty.New()
	_, err := client.R().SetResult(&searchJsonData).Post("https://hunter.qianxin.com/openApi/search/batch?" +
		"username=" + userName +
		"&api-key=" + apiKey +
		"&search=" + base64UrlEncode(search) +
		"&start_time=" + startTime +
		"&end_time=" + endTime)
	if searchJsonData.Message != "success" {
		log.Fatalln("searchApi调用错误提示：" + searchJsonData.Message)
		return searchJsonData, err
	}
	c := 1
check:
	for true {
		_, err = client.R().SetResult(&searchJsonData).Get("https://hunter.qianxin.com/openApi/search/batch/" +
			strconv.Itoa(searchJsonData.Data.Task_id) +
			"?username=" + userName +
			"&api-key=" + apiKey)
		if searchJsonData.Message != "success" {
			log.Fatalln("批量查询任务进度查询接口调用错误提示：" + searchJsonData.Message)
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
	if err != nil {
		log.Fatalf("文件下载接口调用失败 #%v", err)
		return err
	}
	if resp.Body == nil {
		log.Fatalf("下载文件内容为空 #%v", err)
		return err
	}
	defer resp.Body.Close()
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("创建下载文件失败 #%v", err)
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		log.Fatalf("下载文件写入本地失败 #%v", err)
		return err
	}
	return nil
}
