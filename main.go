package main

import (
	"HunterX/obj"
	"HunterX/util"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	page          int
	pageSize      int
	query         string
	searchAll     bool
	batchFilePath string
	outPath       string
	qyLine        bool
	startTime     string
	endTime       string
	isInc         bool
	isHB          bool
	apiUrl        = "https://hunter.qianxin.com"
)

//读取命令行输入
func getFlag() {
	flag.StringVar(&batchFilePath, "l", "", "批量语法查询全部,查询语法文件txt位置")
	flag.BoolVar(&searchAll, "all", false, "查询所有结果")
	flag.IntVar(&page, "page", 1, "单语法查询分页：页数。默认：1")
	flag.IntVar(&pageSize, "size", 10, "单语法查询分页：每页条数。默认：10")
	flag.StringVar(&query, "q", "", "查询语句")
	flag.BoolVar(&qyLine, "qyLine", false, "使用企业积分进行查询")
	flag.StringVar(&endTime, "eTime", time.Now().Format("2006-01-02"), "结束时间，默认当前时间")
	flag.StringVar(&startTime, "sTime", time.Now().AddDate(-1, 0, 0).Format("2006-01-02"), "开始时间，默认一年前当前时间")
	flag.BoolVar(&isInc, "inc", false, "使用内网接口")
	flag.BoolVar(&isHB, "hb", false, "批量查询结果是否合并为同一文件。合并：true，默认不合并：false")
	//加载命令行输入参数
	flag.Parse()
}

func main() {
	println(" _   _             _          __  __\n| | | |_   _ _ __ | |_ ___ _ _\\ \\/ /\n| |_| | | | | '_ \\| __/ _ \\ '__\\  / \n|  _  | |_| | | | | ||  __/ |  /  \\ \n|_| |_|\\__,_|_| |_|\\__\\___|_| /_/\\_\\\n                           v1.3 from:7estUser\n ")
	getFlag()
	if query == "" && batchFilePath == "" {
		log.Fatalf("单语法查询 -q 参数为必须参数，不能为空")
	}
	optionFile, err := ioutil.ReadFile("./hunterx.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败 #%v", err)
	}
	//读取用户查询密钥
	var hunterxConfig obj.Config
	err = yaml.Unmarshal(optionFile, &hunterxConfig)
	if err != nil {
		log.Fatalf("配置文件解析失败 #%v", err)
	}
	userName := hunterxConfig.UserName
	apiKey := hunterxConfig.ApiKey
	//使用内网接口
	if isInc {
		apiUrl = "https://inner.hunter.qianxin-inc.cn"
	}
	log.Println("Hunter查询中，请勿关闭进程...")
	//单语法查询
	if batchFilePath == "" {
		//全部查询
		if searchAll {
			//使用个人账号
			if !qyLine {
				//创建导出Excel
				outFile := xlsx.NewFile()
				//创建结果导出文件随即文件名
				outFileName := util.OutFileName(query)
				util.InitExcel(outFile)
				defer outFile.Save(outFileName)
				///分页遍历查询所有
				searchErr := searchAllDataFor(userName, apiKey, query, startTime, endTime, outFile)
				log.Print("结果保存到文件：" + outFileName)
				if searchErr != nil {
					log.Println("searchApi调用失败: #%v", searchErr)
					return
				}

			} else /*使用企业账号*/ {
				//通过导出接口下载所有数据
				searchAllData(userName, apiKey, query, startTime, endTime)
			}
		} else /*分页查询*/ {
			searchResultData, searchErr := searchData(userName, apiKey, query, page, pageSize, startTime, endTime)
			if searchErr == nil {
				//创建导出Excel
				outFile := xlsx.NewFile()
				//创建结果导出文件随即文件名
				outFileName := util.OutFileName(query)
				util.InitExcel(outFile)
				defer outFile.Save(outFileName)
				//查询结果导出
				util.WriteExcelNew(outFile, searchResultData)
				log.Println("查询完成！结果总数量：" + strconv.Itoa(searchResultData.Data.Total) + "; " + searchResultData.Data.Consume_quota + "; " + searchResultData.Data.Rest_quota + "; 结果保存到文件：" + outFileName)
			}
		}
	} else /*批量查询*/ {
		_, err := os.Stat(batchFilePath)
		if err == nil {
			targetfile, openfIleErr := os.Open(batchFilePath)
			defer targetfile.Close()
			if openfIleErr != nil {
				log.Fatalf("批量查询目标文件打开失败 #%v", openfIleErr)
			}
			scanner := bufio.NewScanner(targetfile)
			if isHB {
				//创建导出Excel
				outFile := xlsx.NewFile()
				util.InitExcel(outFile)
				defer outFile.Save("批量查询汇总.xlsx")
				for scanner.Scan() {
					println("开始查询：" + scanner.Text())
					if !qyLine {
						searchAllDataFor(userName, apiKey, scanner.Text(), startTime, endTime, outFile)
					} else {
						searchAllData(userName, apiKey, scanner.Text(), startTime, endTime)
					}
					time.Sleep(time.Second * 3)
				}

			} else {
				for scanner.Scan() {
					println("开始查询：" + scanner.Text())
					//创建导出Excel
					outFile := xlsx.NewFile()
					//创建结果导出文件随即文件名
					outFileName := util.OutFileName(scanner.Text())
					util.InitExcel(outFile)
					defer outFile.Save(outFileName)
					if !qyLine {
						searchAllDataFor(userName, apiKey, scanner.Text(), startTime, endTime, outFile)
					} else {
						searchAllData(userName, apiKey, scanner.Text(), startTime, endTime)
					}
					time.Sleep(time.Second * 3)
				}
			}
		} else {
			log.Fatalf("批量查询语法文件 " + batchFilePath + " 不存在，请检查")
		}
	}
}

//分页遍历查询所有数据并导出
func searchAllDataFor(userName string, apiKey string, search string, start_time string, end_time string, outFile *xlsx.File) error {
	//通过查询一条数据获取本次查询数据总数量
	searchData, err := util.SearchApi(apiUrl, userName, apiKey, search, 1, 1, start_time, end_time)
	if err != nil {
		log.Println("searchApi调用失败 #%v", err)
		return err
	}
	if searchData.Code == 200 && strings.EqualFold("success", searchData.Message) {
		//查询数据总数量
		title := searchData.Data.Total
		if title > 10000 {
			fmt.Printf("本次查询总条数为：%d，超过1w条限制，仅查询前10000条数据\n", title)
			title = 10000
		}
		//计算每页100条总页数
		pageMax := (title-1)/100 + 1
		//每页100条进行遍历查询
		for j := 1; j <= pageMax; j++ {
			time.Sleep(time.Second * 2)
			searchJsonData, _ := util.SearchApi(apiUrl, userName, apiKey, search, j, 100, start_time, end_time)
			if err != nil {
				return err
			}
			//导出结果
			util.WriteExcelNew(outFile, searchJsonData)
		}
		log.Print(query + "查询完成！结果总数量：" + strconv.Itoa(searchData.Data.Total) + "; " + searchData.Data.Consume_quota + "; " + searchData.Data.Rest_quota + "; ")
		return nil
	} else {
		return errors.New(searchData.Message)
	}

}

//分页查询指定条数并导出结果
func searchData(userName string, apiKey string, search string, p int, s int, start_time string, end_time string) (searchJsonData obj.SearchObj, err error) {
	//分页查询
	searchResultData, err := util.SearchApi(apiUrl, userName, apiKey, search, p, s, start_time, end_time)
	if err != nil {
		log.Println("searchApi调用失败 #%v", err)
		return searchJsonData, err
	}
	return searchResultData, nil
}

//通过接口查询所有并导出结果
func searchAllData(userName string, apiKey string, search string, start_time string, end_time string) {
	searchData, err := util.SearchAllApi(apiUrl, userName, apiKey, search, start_time, end_time)
	if err != nil {
		log.Fatalf("批量查询接口调用失败 #%v", err)
	}
	if strconv.Itoa(searchData.Data.Task_id) != "" && searchData.Data.Progress == "100%" {
		downloadUrl := "https://inner.hunter.qianxin-inc.cn/openApi/search/download/" + strconv.Itoa(searchData.Data.Task_id) +
			"?username=" + userName +
			"&api-key=" + apiKey
		outFileName := util.OutFileName(search)
		err = util.DownloadFile(downloadUrl, outFileName+".csv")
		if err != nil {
			log.Fatalf("全部查询结果文件下载失败 #%v", err)
		}
		println("查询完成！结果总数量：" + strconv.Itoa(searchData.Data.Total) + "; " + searchData.Data.Consume_quota + "; " + searchData.Data.Rest_quota + "; 结果保存到文件：" + outFileName)
	}
}
