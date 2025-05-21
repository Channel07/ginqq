package ginqq

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

// ChainBuilder 链式构建器，用于构建中间件链，
type ChainBuilder struct {
	middlewares []http.RoundTripper // 中间件列表
	head        http.RoundTripper   // 调用链头指针
}

func NewChainBuilder(base http.RoundTripper) *ChainBuilder {
	if base == nil {
		panic("base transport cannot be nil")
	}
	return &ChainBuilder{head: base}
}

// Use 构建中间件链表，头插法
func (b *ChainBuilder) Use(middlewares ...http.RoundTripper) *ChainBuilder {
	for i := len(middlewares) - 1; i >= 0; i-- { // 反向遍历，数组中第一个元素作为头指针，
		m := middlewares[i]
		if injector, ok := m.(interface{ SetNext(http.RoundTripper) }); ok {
			injector.SetNext(b.head)
			b.head = m
		} else {
			panic("middleware must implement SetNext interface")
		}
	}
	b.middlewares = append(b.middlewares, middlewares...)
	return b
}

// Build 返回调用链头指针
func (b *ChainBuilder) Build() http.RoundTripper {
	return b.head
}

// LoggingTripper2 中间件示例，需要实现RoundTrip处理接口，SetNext设置下个中间接口
type LoggingTripper2 struct {
	next http.RoundTripper
}

func NewLoggingTripper2() *LoggingTripper2 {
	return &LoggingTripper2{}
}

func (l *LoggingTripper2) SetNext(next http.RoundTripper) {
	l.next = next
}

func (l *LoggingTripper2) RoundTrip(req *http.Request) (*http.Response, error) {

	log.Printf("[TRIPPER2] BEFORE %s %s", req.Method, req.URL)

	start := time.Now()
	resp, err := l.next.RoundTrip(req) // 调用下一层

	log.Printf(
		"[TRIPPER2] AFTER %s %d %v",
		req.URL.Path, resp.StatusCode, time.Since(start),
	)
	return resp, err
}

type LoggingTripper struct {
	next http.RoundTripper
}

func NewLoggingTripper() *LoggingTripper {
	return &LoggingTripper{}
}

func (l *LoggingTripper) SetNext(next http.RoundTripper) {
	l.next = next
}

func (l *LoggingTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	log.Printf("[TRIPPER] BEFORE %s %s", req.Method, req.URL)

	start := time.Now()
	resp, err := l.next.RoundTrip(req) // 调用下一层

	log.Printf(
		"[TRIPPER] AFTER %s %d %v",
		req.URL.Path, resp.StatusCode, time.Since(start),
	)
	return resp, err
}

func HttpEnhance(cfg *HttpClientEnhanceConfig) {

	base := cfg.Transport
	if base == nil {
		base = http.DefaultTransport
	}
	// 跳过证书认证
	if !cfg.DisableSkipVerify {
		if base.(*http.Transport).TLSClientConfig == nil {
			base.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			base.(*http.Transport).TLSClientConfig.InsecureSkipVerify = true
		}
	}

	// 注册中间件
	var middlewares []http.RoundTripper
	if true {
		middlewares = append(middlewares, NewLoggingTripper())
	}
	if true {
		middlewares = append(middlewares, NewLoggingTripper2())
	}

	// 强制替换默认传输层为增强的传输层
	http.DefaultTransport = NewChainBuilder(base).
		Use(middlewares...).
		Build()
}
