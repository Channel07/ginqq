package ginqq

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
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

// DispatchTransactionLog 调度内部流水日志，作为中间件使用。
func DispatchTransactionLog(c *Context) {
	// 1.创建流水日志对象
	log := &TransactionLog{ctx: c}

	// 2.在处理请求之前填充的日志字段
	log.before()

	// 3.处理请求
	requestTime := time.Now()
	c.Next()
	responseTime := time.Now()

	// 4.在处理请求之后填充的日志字段，并在填充后记录到文件
	go func() {
		log.after(requestTime, responseTime)
		log.logger()
	}()

	fmt.Println("over")
}

// before 在处理请求 之前 填充的字段（有个别字段是必须在处理请求之前获取的，放到 before 中执行）。
func (log *TransactionLog) before() {
	defer deferRecover()
	log.GetRequestPayload()
}

// after 在处理请求 之后 填充的字段。
func (log *TransactionLog) after(requestTime, responseTime time.Time) {
	defer deferRecover()
	log.GetAppName()
	log.GetLevel()
	log.GetLogTime()
	log.GetLogger()
	log.GetThread()
	log.GetTransactionID()
	log.GetDialogType()
	log.GetAddress()
	log.GetFCode()
	log.GetTCode()
	log.GetMethodCode()
	log.GetMethodName()
	log.GetHTTPMethod()
	log.GetRequestTime(requestTime)
	log.GetRequestHeaders()
	log.GetResponseTime(responseTime)
	log.GetResponseHeaders()
	log.GetResponsePayload()
	log.GetResponseRemark()
	log.GetResponseCode()
	log.GetHTTPStatusCode()
	log.GetOrderID()
	log.GetProvinceCode()
	log.GetCityCode()
	log.GetTotalTime(requestTime, responseTime)
	log.GetErrorCode()
	log.GetRequestIP()
	log.GetHostIP()
	log.GetHostname()
	log.GetAccountType()
	log.GetAccountNum()
	log.GetResponseAccountType()
	log.GetResponseAccountNum()
	log.GetUser()
	log.GetTag()
	log.GetServiceLine()
}

// logger 记录日志到文件
func (log *TransactionLog) logger() {
	defer deferRecover()
	logJson, err := json.Marshal(log)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(logJson))
}

func (log *TransactionLog) GetAppName() *TransactionLog {
	log.AppName = strings.ToLower(cnf.SvcCode) + "_" + cnf.AppName
	return log
}

func (log *TransactionLog) GetLevel() *TransactionLog {
	log.Level = "INFO"
	return log
}

func (log *TransactionLog) GetLogTime() *TransactionLog {
	log.LogTime = time.Now().Format("2006-01-02 15:04:05.000")
	return log
}

func (log *TransactionLog) GetLogger() *TransactionLog {
	log.Logger = "ginqq"
	return log
}

// GetThread 获取当前线程ID（这里实际获取的是 Goroutine ID）。
func (log *TransactionLog) GetThread() *TransactionLog {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	goroutineID := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	log.Thread = goroutineID
	return log
}

func (log *TransactionLog) GetTransactionID() *TransactionLog {
	log.TransactionID = log.ctx.GetTransactionID()
	return log
}

func (log *TransactionLog) GetDialogType() *TransactionLog {
	log.DialogType = "in"
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

func (log *TransactionLog) GetFCode() *TransactionLog {
	log.FCode = log.ctx.GetFCode()
	return log
}

func (log *TransactionLog) GetTCode() *TransactionLog {
	log.TCode = cnf.SvcCode
	return log
}

func (log *TransactionLog) GetMethodCode() *TransactionLog {
	log.MethodCode = log.ctx.GetMethodCode()
	return log
}

func (log *TransactionLog) GetMethodName() *TransactionLog {
	log.MethodName = log.ctx.GetMethodName()
	return log
}

func (log *TransactionLog) GetHTTPMethod() *TransactionLog {
	log.HTTPMethod = log.ctx.Request.Method
	return log
}

func (log *TransactionLog) GetRequestTime(requestTime time.Time) *TransactionLog {
	log.RequestTime = requestTime.Format("2006-01-02 15:04:05.000")
	return log
}

func (log *TransactionLog) GetRequestHeaders() *TransactionLog {
	headers := make(map[string]string)
	for k, v := range log.ctx.Request.Header {
		headers[k] = strings.Join(v, ", ")
	}
	headersSerialized, _ := json.Marshal(headers)
	log.RequestHeaders = string(headersSerialized)
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

func (log *TransactionLog) GetResponseTime(responseTime time.Time) *TransactionLog {
	log.ResponseTime = responseTime.Format("2006-01-02 15:04:05.000")
	return log
}

func (log *TransactionLog) GetResponseHeaders() *TransactionLog {
	headers := make(map[string]string)
	for k, v := range log.ctx.Writer.Header() {
		headers[k] = strings.Join(v, ", ")
	}
	headersSerialized, _ := json.Marshal(headers)
	log.ResponseHeaders = string(headersSerialized)
	return log
}

func (log *TransactionLog) GetResponsePayload() *TransactionLog {
	responsePayload := log.ctx.GetResponsePayload()
	if responsePayload != nil {
		responsePayloadSerialized, _ := json.Marshal(responsePayload)
		log.ResponsePayload = string(responsePayloadSerialized)
	} else {
		log.ResponsePayload = "{}"
	}
	return log
}

func (log *TransactionLog) GetResponseRemark() *TransactionLog {
	return log
}

func (log *TransactionLog) GetResponseCode() *TransactionLog {
	responsePayload := log.ctx.GetResponsePayload()

	if responsePayload != nil {
		// TODO if responsePayload is not map.
		var deserialization interface{}
		serialized, _ := json.Marshal(responsePayload)
		_ = json.Unmarshal(serialized, &deserialization)
		code := FuzzyGet(deserialization, "code")
		if code != nil {
			log.ResponseCode = fmt.Sprintf("%v", code)
		}
	}

	return log
}

func (log *TransactionLog) GetHTTPStatusCode() *TransactionLog {
	log.HTTPStatusCode = strconv.Itoa(log.ctx.Writer.Status())
	return log
}

func (log *TransactionLog) GetOrderID() *TransactionLog {
	return log
}

func (log *TransactionLog) GetProvinceCode() *TransactionLog {
	return log
}

func (log *TransactionLog) GetCityCode() *TransactionLog {
	return log
}

func (log *TransactionLog) GetTotalTime(requestTime, responseTime time.Time) *TransactionLog {
	log.TotalTime = responseTime.Sub(requestTime).Milliseconds()
	return log
}

func (log *TransactionLog) GetErrorCode() *TransactionLog {
	return log
}

func (log *TransactionLog) GetRequestIP() *TransactionLog {
	return log
}

func (log *TransactionLog) GetHostIP() *TransactionLog {
	return log
}

func (log *TransactionLog) GetHostname() *TransactionLog {
	return log
}

func (log *TransactionLog) GetAccountType() *TransactionLog {
	return log
}

func (log *TransactionLog) GetAccountNum() *TransactionLog {
	return log
}

func (log *TransactionLog) GetResponseAccountType() *TransactionLog {
	return log
}

func (log *TransactionLog) GetResponseAccountNum() *TransactionLog {
	return log
}

func (log *TransactionLog) GetUser() *TransactionLog {
	return log
}

func (log *TransactionLog) GetTag() *TransactionLog {
	return log
}

func (log *TransactionLog) GetServiceLine() *TransactionLog {
	return log
}

func deferRecover() {
	if err := recover(); err != nil {
		fmt.Printf("An error occurred while executing the transaction log middleware：%v\n", err)
	}
}
