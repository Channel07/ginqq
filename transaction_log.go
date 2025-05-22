package ginqq

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TransactionLog struct {
	ctx                 *Context
	AppName             string `json:"app_name"`
	Level               string `json:"level"`
	LogTime             string `json:"log_time"`
	Logger              string `json:"logger"`
	Thread              string `json:"thread"`
	TransactionID       string `json:"transaction_id"`
	DialogType          string `json:"dialog_type"`
	Address             string `json:"address"`
	FCode               string `json:"fcode"`
	TCode               string `json:"tcode"`
	MethodCode          string `json:"method_code"`
	MethodName          string `json:"method_name"`
	HTTPMethod          string `json:"http_method"`
	RequestTime         string `json:"request_time"`
	RequestHeaders      string `json:"request_headers"`
	RequestPayload      string `json:"request_payload"`
	ResponseTime        string `json:"response_time"`
	ResponseHeaders     string `json:"response_headers"`
	ResponsePayload     string `json:"response_payload"`
	ResponseCode        string `json:"response_code"`
	ResponseRemark      string `json:"response_remark"`
	HTTPStatusCode      string `json:"http_status_code"`
	OrderID             string `json:"order_id"`
	ProvinceCode        string `json:"province_code"`
	CityCode            string `json:"city_code"`
	TotalTime           int64  `json:"total_time"`
	ErrorCode           string `json:"error_code"`
	RequestIP           string `json:"request_ip"`
	HostIP              string `json:"host_ip"`
	Hostname            string `json:"hostname"`
	AccountType         string `json:"account_type"`
	AccountNum          string `json:"account_num"`
	ResponseAccountType string `json:"response_account_type"`
	ResponseAccountNum  string `json:"response_account_num"`
	User                string `json:"user"`
	Tag                 string `json:"tag"`
	ServiceLine         string `json:"service_line"`
}

func DispatchTransactionLog(c *Context) {
	log := &TransactionLog{ctx: c}

	before(log)

	requestTime := time.Now()
	c.Next()
	responseTime := time.Now()

	after(log, requestTime, responseTime)

	logJson, err := json.Marshal(log)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(logJson))
}

func before(log *TransactionLog) {
	defer deferRecover()
	log.GetRequestPayload()
}

func after(log *TransactionLog, requestTime, responseTime time.Time) {
	defer deferRecover()
	go log.GetAppName()
	go log.GetAddress()
	// 该方法链中的两个方法会在同一个新的goroutine中顺序执行：
	//go log.GetAppName().GetAddress()
}

func (log *TransactionLog) GetAppName() *TransactionLog {
	svcCode := strings.ToLower(globalConfig.SvcCode)
	appName := strings.ReplaceAll(strings.ToLower(globalConfig.AppName), "-", "_")
	log.AppName = svcCode + "_" + appName
	return log
}

func (log *TransactionLog) GetAddress() *TransactionLog {
	var schema string
	if log.ctx.Request.TLS != nil {
		schema = "https"
	} else {
		schema = "http"
	}
	host := log.ctx.Request.Host
	path := log.ctx.Request.URL.Path
	log.Address = fmt.Sprintf("%s//%s%s", schema, host, path)
	return log
}

// GetRequestPayload 获取请求数据，仅获取查询参数，表单数据，JSON数据。
func (log *TransactionLog) GetRequestPayload() *TransactionLog {
	requestPayload := make(map[string]interface{})

	// Get query params and form data.
	_ = log.ctx.Request.ParseForm()
	for key, values := range log.ctx.Request.Form {
		if len(values) == 1 {
			requestPayload[key] = values[0]
		} else {
			requestPayload[key] = values
		}
	}

	// Get json data.
	requestBody, _ := log.ctx.GetRawDataReusable()
	if len(requestBody) != 0 {
		_ = json.Unmarshal(requestBody, &requestPayload)
	}

	requestPayloadSerialized, _ := json.Marshal(requestPayload)
	log.RequestPayload = string(requestPayloadSerialized)

	return log
}

func deferRecover() {
	if err := recover(); err != nil {
		fmt.Printf("An error occurred while executing the transaction log middleware：%v\n", err)
	}
}
