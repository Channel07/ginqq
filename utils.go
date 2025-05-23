package ginqq

import (
	uuid "github.com/satori/go.uuid"
	"reflect"
	"runtime"
	"strings"
)

func uuid4() string {
	u4 := uuid.NewV4()
	return strings.ReplaceAll(u4.String(), "-", "")
}

// getRawHandlerName 获取原始处理函数名称
func getRawHandlerName(h func(*Context)) string {
	// 获取函数指针
	ptr := reflect.ValueOf(h).Pointer()

	// 通过 runtime 获取完整名称
	fullName := runtime.FuncForPC(ptr).Name()

	// 处理结构体方法名称（如：github.com/xxx.(*Type).MethodName）
	if strings.Contains(fullName, ").") {
		parts := strings.Split(fullName, ").")
		if len(parts) > 1 {
			// 示例：将 "github.com/xxx.(*Type).MethodName-fm" → "MethodName"
			name := strings.Split(parts[1], "-")[0]
			return strings.TrimSuffix(name, ".fm")
		}
	}

	// 处理闭包名称（如：convertToGinHandlers.func1）
	if strings.Contains(fullName, ".func") {
		parts := strings.Split(fullName, ".")
		// 取闭包外层函数名（示例：convertToGinHandlers）
		if len(parts) >= 2 {
			return parts[len(parts)-2]
		}
	}

	// 通用处理：取最后一段作为名称
	parts := strings.Split(fullName, ".")
	return parts[len(parts)-1]
}

// FuzzyGet 在嵌套数据结构中模糊查找键并返回值
func FuzzyGet(data interface{}, key string) interface{} {
	return fuzzyGet(data, processKey(key))
}

// FuzzyGet 辅助函数，递归处理数据结构
func fuzzyGet(data interface{}, targetKey string) interface{} {
	// 处理 map 类型
	if m, ok := data.(map[string]interface{}); ok {
		// 先检查当前层
		for k, v := range m {
			if processKey(k) == targetKey {
				return v
			}
		}
		// 递归处理子节点
		for _, v := range m {
			if result := fuzzyGet(v, targetKey); result != nil {
				return result
			}
		}
		return nil
	}

	// 使用反射处理所有切片/数组类型
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i).Interface()
			if result := fuzzyGet(elem, targetKey); result != nil {
				return result
			}
		}
	}

	return nil
}

// FuzzyGet 辅助函数，处理键格式。
func processKey(key string) string {
	key = strings.ReplaceAll(key, " ", "")
	key = strings.ReplaceAll(key, "-", "")
	key = strings.ReplaceAll(key, "_", "")
	return strings.ToLower(key)
}
