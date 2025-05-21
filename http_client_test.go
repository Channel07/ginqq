package ginqq

import (
	"net/http"
	"testing"
)

func TestHttpClient(t *testing.T) {
	Default("Y122010101", "Y122")
	http.Get("http://www.baidu.com")
}
