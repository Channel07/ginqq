package ginqq

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TransactionLog struct {
	ctx *Context

	requestTime  time.Time
	responseTime time.Time

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
	ResponseRemark      string `json:"response_remark"`
	ResponseCode        string `json:"response_code"`
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
	log := &TransactionLog{ctx: c}

	log.before()

	log.requestTime = time.Now()
	c.Next()
	log.responseTime = time.Now()

	go func() {
		log.after()
		msg, _ := json.Marshal(log)
		logger.Info(string(msg))
	}()
}

func (log *TransactionLog) before() {
	defer deferRecover()
	log.GetRequestPayload()
}

func (log *TransactionLog) after() {
	defer deferRecover()

	var wg sync.WaitGroup

	v := reflect.ValueOf(log)
	typ := v.Type()

	for i := 0; i < typ.NumMethod(); i++ {
		methodName := typ.Method(i).Name
		if strings.HasPrefix(methodName, "Get") && methodName != "GetRequestPayload" {
			wg.Add(1)
			go func(m reflect.Value, name string) {
				defer func() {
					wg.Done()
					if err := recover(); err != nil {
						fmt.Printf("[TransactionLog] %s panic: %v\n", name, err)
						debug.PrintStack()
					}
				}()
				m.Call(nil)
			}(v.Method(i), methodName)
		}
	}

	wg.Wait()
}

func (log *TransactionLog) GetAppName() *TransactionLog {
	log.AppName = strings.ToLower(cnf.SvcCode) + "_" + cnf.AppName + "_info"
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

func (log *TransactionLog) GetRequestTime() *TransactionLog {
	log.RequestTime = log.requestTime.Format("2006-01-02 15:04:05.000")
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

func (log *TransactionLog) GetResponseTime() *TransactionLog {
	log.ResponseTime = log.responseTime.Format("2006-01-02 15:04:05.000")
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
	if responsePayload := log.ctx.GetResponsePayload(); responsePayload != nil {
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
	if responsePayload := log.ctx.GetResponsePayload(); responsePayload != nil {
		log.ResponseCode = FuzzyGet(CrossJson(responsePayload), "code")
	}
	return log
}

func (log *TransactionLog) GetHTTPStatusCode() *TransactionLog {
	log.HTTPStatusCode = strconv.Itoa(log.ctx.Writer.Status())
	return log
}

func (log *TransactionLog) GetOrderID() *TransactionLog {
	if responsePayload := log.ctx.GetResponsePayload(); responsePayload != nil {
		log.OrderID = FuzzyGetMany(CrossJson(responsePayload), []string{"order_id", "ht_id"})
	}
	return log
}

func (log *TransactionLog) GetProvinceCodeAndCityCode() *TransactionLog {
	var requestPayload map[string]interface{}
	_ = json.Unmarshal([]byte(log.RequestPayload), &requestPayload)
	responsePayload := CrossJson(log.ctx.GetResponsePayload())
	mergedPayload := []interface{}{requestPayload, responsePayload}
	log.ProvinceCode = FuzzyGet(mergedPayload, "province_code")
	log.CityCode = FuzzyGet(mergedPayload, "city_code")
	return log
}

func (log *TransactionLog) GetTotalTime() *TransactionLog {
	log.TotalTime = log.responseTime.Sub(log.requestTime).Milliseconds()
	return log
}

func (log *TransactionLog) GetErrorCode() *TransactionLog {
	return log
}

func (log *TransactionLog) GetRequestIP() *TransactionLog {
	log.RequestIP = log.ctx.ClientIP()
	return log
}

func (log *TransactionLog) GetHostIP() *TransactionLog {
	if hostIP, err := GetHostIP(); err == nil {
		log.HostIP = hostIP
	}
	return log
}

func (log *TransactionLog) GetHostname() *TransactionLog {
	if hostname, err := os.Hostname(); err == nil {
		log.Hostname = hostname
	}
	return log
}

func (log *TransactionLog) GetAccount() *TransactionLog {
	var requestPayload map[string]interface{}
	if err := json.Unmarshal([]byte(log.RequestPayload), &requestPayload); err == nil {
		keys := []string{"phone", "phone_num", "number", "accnbr"}
		accountNum := FuzzyGetMany(requestPayload, keys)
		if accountNum != "" {
			log.AccountType = "11"
			log.AccountNum = accountNum
		}
	}
	return log
}

func (log *TransactionLog) GetResponseAccount() *TransactionLog {
	if responsePayload := log.ctx.GetResponsePayload(); responsePayload != nil {
		keys := []string{"phone", "phone_num", "accnbr", "receive_phone"}
		responseAccountNum := FuzzyGetMany(CrossJson(responsePayload), keys)
		if responseAccountNum != "" {
			log.ResponseAccountType = "11"
			log.ResponseAccountNum = responseAccountNum
		}
	}
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
		debug.PrintStack()
	}
}
