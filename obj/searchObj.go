package obj

type SearchObj struct {
	//成功：200
	Code int `json:"code"`
	//成功：success
	Message string `json:"message"`
	//查询数据
	Data searchData `json:"data"`
}

type searchData struct {
	//消耗积分
	Consume_quota string `json:"consume_quota"`
	//今日剩余积分
	Rest_quota string `json:"rest_quota"`
	//查询结果总数量
	Total int `json:"total"`
	//查询消耗时间
	Time int `json:"time"`
	//查询结果数据
	Arr []arrData `json:"arr"`
	//账号类型
	Account_type string `json:"account_type"`
	//批量查询：任务ID
	Task_id int `json:"task_id"`
	//批量查询：导出文件名称
	FileName string `json:"filename"`
	//批量查询：任务执行进度
	Progress string `json:"progress"`
}

type arrData struct {
	Is_risk          string `json:"is_risk"`
	Url              string `json:"url"`
	Ip               string `json:"ip"`
	Port             int    `json:"port"`
	Web_title        string `json:"web_title"`
	Domain           string `json:"domain"`
	Is_risk_protocol string `json:"is_risk_protocol"`
	//协议
	Protocol string `json:"protocol"`
	//通信协议
	Base_protocol string `json:"base_protocol"`
	//网站状态码
	Status_code int             `json:"status_code"`
	Component   []componentData `json:"component"`
	//操作系统
	Os         string `json:"os"`
	Country    string `json:"country"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Updated_at string `json:"updated_at"`
	//备案单位
	Company string `json:"company"`
	Number  string `json:"number"`
	Is_web  string `json:"is_web"`
	As_org  string `json:"as_org"`
	Isp     string `json:"isp"`
	Banner  string `json:"banner"`
}

type componentData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
