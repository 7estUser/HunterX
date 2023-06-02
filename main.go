package main

import (
	"HunterX/obj"
	"HunterX/util"
	"bufio"
	"flag"
	"github.com/xuri/excelize/v2"
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
)

//读取命令行输入
func getFlag() {
	flag.StringVar(&batchFilePath, "l", "", "批量语法查询全部,查询语法文件txt位置")
	flag.BoolVar(&searchAll, "all", false, "是否查询所有结果,默认为false,设置true时分页失效,查询所有结果")
	flag.IntVar(&page, "page", 1, "单语法查询分页：页数。默认：1")
	flag.IntVar(&pageSize, "size", 10, "单语法查询分页：每页条数。默认：10")
	flag.StringVar(&query, "q", "", "单语法查询语句")
	flag.BoolVar(&qyLine, "qyLine", false, "使用企业导出额度进行全部查询，只针对企业账号，默认为false：使用权益积分进行查询")
	flag.StringVar(&endTime, "eTime", time.Now().Format("2006-01-02"), "结束时间，默认当前时间")
	flag.StringVar(&startTime, "sTime", time.Now().AddDate(-1, 0, 0).Format("2006-01-02"), "开始时间，默认一年前当前时间")
	//加载命令行输入参数
	flag.Parse()
}

func main() {
	println(" _   _             _          __  __\n| | | |_   _ _ __ | |_ ___ _ _\\ \\/ /\n| |_| | | | | '_ \\| __/ _ \\ '__\\  / \n|  _  | |_| | | | | ||  __/ |  /  \\ \n|_| |_|\\__,_|_| |_|\\__\\___|_| /_/\\_\\\n                           v1.0 from:duz\n ")
	getFlag()
	if query == "" && batchFilePath == "" {
		println("单语法查询 -q 参数为必须参数，不能为空")
		return
	}
	optionFile, err := ioutil.ReadFile("./hunterx.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败 #%v", err)
		return
	}
	var hunterxConfig obj.Config
	err = yaml.Unmarshal(optionFile, &hunterxConfig)
	if err != nil {
		log.Fatalf("配置文件解析失败 #%v", err)
		return
	}
	userName := hunterxConfig.UserName
	apiKey := hunterxConfig.ApiKey
	println("Hunter查询中，请勿关闭进程...")
	//单语法查询
	if batchFilePath == "" {
		//全部查询,
		if searchAll {
			//查询用户类型
			//userAccountType, _ := util.SelectAccountType(userName, apiKey)
			//if userAccountType != "个人账号" && !qyLine {
			if !qyLine {
				searchAllDataFor(userName, apiKey, query, startTime, endTime)
			} else {
				searchAllData(userName, apiKey, query, startTime, endTime)
			}

		} else {
			searchData(userName, apiKey, query, page, pageSize, startTime, endTime)
		}
	} else {
		_, err := os.Stat(batchFilePath)
		if err == nil {
			targetfile, openfIleErr := os.Open(batchFilePath)
			defer targetfile.Close()
			if openfIleErr != nil {
				log.Fatalf("批量查询目标文件打开失败 #%v", openfIleErr)
				return
			}
			scanner := bufio.NewScanner(targetfile)
			for scanner.Scan() {
				println(scanner.Text())
				//if userAccountType != "个人账号" && !qyLine {
				if !qyLine {
					searchAllDataFor(userName, apiKey, scanner.Text(), startTime, endTime)
				} else {
					searchAllData(userName, apiKey, scanner.Text(), startTime, endTime)
				}
				time.Sleep(time.Second * 3)
			}
		} else {
			log.Fatalf("批量查询语法文件 " + batchFilePath + " 不存在，请检查")
		}
	}
}

//使用企业导出额度查询所有结果并导出
func searchAllDataFor(u string, k string, search string, start_time string, end_time string) {
	searchData, err := util.SearchApi(u, k, search, 1, 1, start_time, end_time)
	if err != nil {
		log.Fatalf("searchApi调用失败 #%v", err)
		return
	}
	if searchData.Code == 200 && strings.EqualFold("success", searchData.Message) {
		//导出结果到excel
		outFile := excelize.NewFile()
		defer outFile.Close()
		util.InitExcel(outFile)
		//查询数据总数量
		title := searchData.Data.Total
		//一共多少页
		pageMax := title/100 + 1
		for j := 1; j <= pageMax; j++ {
			time.Sleep(time.Second * 2)
			searchJsonData, err := util.SearchApi(u, k, search, j, 100, start_time, end_time)
			if err != nil {
				log.Fatalf("searchApi调用失败 #%v", err)
				return
			}
			for i := 0; i < len(searchJsonData.Data.Arr); i++ {
				util.WriteExcel(outFile, j, i, searchJsonData)
			}
		}
		ourFileName := util.OutFileName(search)
		if util.SaveExcel(outFile, ourFileName) != nil {
			log.Fatalf("数据导出到excel失败 #%v", err)
		}
		log.Println("查询完成！结果总数量：" + strconv.Itoa(searchData.Data.Total) + ";" + searchData.Data.Rest_quota + "; 结果保存到文件：" + ourFileName + ".xlsx")
	}
}

//分页查询并导出结果
func searchData(u string, k string, search string, p int, s int, start_time string, end_time string) {
	//分页查询
	searchResultData, err := util.SearchApi(u, k, search, p, s, start_time, end_time)
	if err != nil {
		log.Fatalf("searchApi调用失败 #%v", err)
		return
	}
	outFile := excelize.NewFile()
	defer outFile.Close()
	util.InitExcel(outFile)
	for i := 0; i < len(searchResultData.Data.Arr); i++ {
		util.WriteExcel(outFile, 1, i, searchResultData)
	}
	outFileName := util.OutFileName(search)
	err = util.SaveExcel(outFile, outFileName)
	if err != nil {
		log.Fatalf("查询结果保存到文件失败 #%v", err)
		return
	}
	log.Println("查询完成！结果总数量：" + strconv.Itoa(searchResultData.Data.Total) + "; " + searchResultData.Data.Consume_quota + "; " + searchResultData.Data.Rest_quota + "; 结果保存到文件：" + outFileName + ".xlsx")

}

//查询所有并导出结果
func searchAllData(u string, k string, search string, start_time string, end_time string) {
	searchData, err := util.SearchAllApi(u, k, search, start_time, end_time)
	if err != nil {
		log.Fatalf("批量查询接口调用失败 #%v", err)
		return
	}
	if strconv.Itoa(searchData.Data.Task_id) != "" && searchData.Data.Progress == "100%" {
		downloadUrl := "https://hunter.qianxin.com/openApi/search/download/" + strconv.Itoa(searchData.Data.Task_id) +
			"?username=" + u +
			"&api-key=" + k
		outFileName := util.OutFileName(search)
		err = util.DownloadFile(downloadUrl, outFileName+".csv")
		if err != nil {
			log.Fatalf("全部查询结果文件下载失败 #%v", err)
			return
		}
		log.Println("查询完成！结果总数量：" + strconv.Itoa(searchData.Data.Total) + "; " + searchData.Data.Consume_quota + "; " + searchData.Data.Rest_quota + "; 结果保存到文件：" + outFileName + ".csv")
	}
}
