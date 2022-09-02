package util

import (
	"HunterX/obj"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func SaveExcel(outFile *excelize.File, fileName string) error {
	return outFile.SaveAs(fileName + ".xlsx")
}

func InitExcel(outFile *excelize.File) {
	//格式化excel
	//outFile := excelize.NewFile()
	outFile.SetCellValue("Sheet1", "A1", "url")
	outFile.SetCellValue("Sheet1", "B1", "资产标签")
	outFile.SetCellValue("Sheet1", "C1", "IP")
	outFile.SetCellValue("Sheet1", "D1", "IP标签")
	outFile.SetCellValue("Sheet1", "E1", "端口")
	outFile.SetCellValue("Sheet1", "F1", "网站标题")
	outFile.SetCellValue("Sheet1", "G1", "域名")
	outFile.SetCellValue("Sheet1", "H1", "高危协议")
	outFile.SetCellValue("Sheet1", "I1", "协议")
	outFile.SetCellValue("Sheet1", "J1", "通讯协议")
	outFile.SetCellValue("Sheet1", "K1", "网站状态码")
	outFile.SetCellValue("Sheet1", "L1", "应用/组件")
	outFile.SetCellValue("Sheet1", "M1", "操作系统")
	outFile.SetCellValue("Sheet1", "N1", "备案单位")
	outFile.SetCellValue("Sheet1", "O1", "备案号")
	outFile.SetCellValue("Sheet1", "P1", "国家")
	outFile.SetCellValue("Sheet1", "Q1", "省份")
	outFile.SetCellValue("Sheet1", "R1", "市区")
	outFile.SetCellValue("Sheet1", "S1", "探查时间")
	outFile.SetCellValue("Sheet1", "T1", "Web资产")
	outFile.SetCellValue("Sheet1", "U1", "运营商")
	outFile.SetCellValue("Sheet1", "V1", "注册机构")
}

func WriteExcel(outFile *excelize.File, pageNum int, linNum int, searchJsonData obj.SearchObj) {
	outFile.SetCellValue("Sheet1", "A"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Url)
	outFile.SetCellValue("Sheet1", "B"+strconv.Itoa((pageNum-1)*100+linNum+2), "")
	outFile.SetCellValue("Sheet1", "C"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Ip)
	outFile.SetCellValue("Sheet1", "D"+strconv.Itoa((pageNum-1)*100+linNum+2), "")
	outFile.SetCellValue("Sheet1", "E"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Port)
	outFile.SetCellValue("Sheet1", "F"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Web_title)
	outFile.SetCellValue("Sheet1", "G"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Domain)
	outFile.SetCellValue("Sheet1", "H"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Is_risk_protocol)
	outFile.SetCellValue("Sheet1", "I"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Protocol)
	outFile.SetCellValue("Sheet1", "J"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Base_protocol)
	outFile.SetCellValue("Sheet1", "K"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Status_code)
	outFile.SetCellValue("Sheet1", "L"+strconv.Itoa((pageNum-1)*100+linNum+2), "")
	outFile.SetCellValue("Sheet1", "M"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Os)
	outFile.SetCellValue("Sheet1", "N"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Company)
	outFile.SetCellValue("Sheet1", "O"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Number)
	outFile.SetCellValue("Sheet1", "P"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Country)
	outFile.SetCellValue("Sheet1", "Q"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].City)
	outFile.SetCellValue("Sheet1", "R"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Province)
	outFile.SetCellValue("Sheet1", "S"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Updated_at)
	outFile.SetCellValue("Sheet1", "T"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Is_web)
	outFile.SetCellValue("Sheet1", "U"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].Isp)
	outFile.SetCellValue("Sheet1", "V"+strconv.Itoa((pageNum-1)*100+linNum+2), searchJsonData.Data.Arr[linNum].As_org)
}
