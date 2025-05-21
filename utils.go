package ginqq

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

func GenerateUuid() string {
	id := uuid.NewV4()
	ids := id.String()
	reqId := strings.ReplaceAll(ids, "-", "")
	return reqId
}
