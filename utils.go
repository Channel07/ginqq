package ginqq

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

func uuid4() string {
	u4 := uuid.NewV4()
	return strings.ReplaceAll(u4.String(), "-", "")
}
