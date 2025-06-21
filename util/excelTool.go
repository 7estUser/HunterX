package util

import (
	"HunterX/obj"
	"github.com/tealeg/xlsx"
	"strconv"
)

//格式化excel
func InitExcel(outFile *xlsx.File) {
	sheet, _ := outFile.AddSheet("Sheet1")
	titleRow := sheet.AddRow()
	titleRow.AddCell().Value = "url"
	titleRow.AddCell().Value = "资产标签"
	titleRow.AddCell().Value = "IP"
	titleRow.AddCell().Value = "IP标签"
	titleRow.AddCell().Value = "端口"
	titleRow.AddCell().Value = "网站标题"
	titleRow.AddCell().Value = "域名"
	titleRow.AddCell().Value = "高危协议"
	titleRow.AddCell().Value = "协议"
	titleRow.AddCell().Value = "通讯协议"
	titleRow.AddCell().Value = "网站状态码"
	titleRow.AddCell().Value = "应用/组件"
	titleRow.AddCell().Value = "操作系统"
	titleRow.AddCell().Value = "备案单位"
	titleRow.AddCell().Value = "备案号"
	titleRow.AddCell().Value = "国家"
	titleRow.AddCell().Value = "省份"
	titleRow.AddCell().Value = "市区"
	titleRow.AddCell().Value = "探查时间"
	titleRow.AddCell().Value = "Web资产"
	titleRow.AddCell().Value = "运营商"
	titleRow.AddCell().Value = "注册机构"
}

func WriteExcelNew(outFile *xlsx.File, searchJsonData obj.SearchObj) {
	sheet := outFile.Sheets[0]
	for i := 0; i < len(searchJsonData.Data.Arr); i++ {
		row := sheet.AddRow()
		row.AddCell().Value = searchJsonData.Data.Arr[i].Url
		row.AddCell().Value = ""
		row.AddCell().Value = searchJsonData.Data.Arr[i].Ip
		row.AddCell().Value = ""
		row.AddCell().Value = strconv.Itoa(searchJsonData.Data.Arr[i].Port)
		row.AddCell().Value = searchJsonData.Data.Arr[i].Web_title
		row.AddCell().Value = searchJsonData.Data.Arr[i].Domain
		row.AddCell().Value = searchJsonData.Data.Arr[i].Is_risk_protocol
		row.AddCell().Value = searchJsonData.Data.Arr[i].Protocol
		row.AddCell().Value = searchJsonData.Data.Arr[i].Base_protocol
		row.AddCell().Value = strconv.Itoa(searchJsonData.Data.Arr[i].Status_code)
		row.AddCell().Value = ""
		row.AddCell().Value = searchJsonData.Data.Arr[i].Os
		row.AddCell().Value = searchJsonData.Data.Arr[i].Company
		row.AddCell().Value = searchJsonData.Data.Arr[i].Number
		row.AddCell().Value = searchJsonData.Data.Arr[i].Country
		row.AddCell().Value = searchJsonData.Data.Arr[i].City
		row.AddCell().Value = searchJsonData.Data.Arr[i].Province
		row.AddCell().Value = searchJsonData.Data.Arr[i].Updated_at
		row.AddCell().Value = searchJsonData.Data.Arr[i].Is_web
		row.AddCell().Value = searchJsonData.Data.Arr[i].Isp
		row.AddCell().Value = searchJsonData.Data.Arr[i].As_org
	}
}
