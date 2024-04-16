package helpers

import (
	"fmt"
	"runtime"
)

type Common struct {
}

func (c *Common) GetErrorLocation() map[string]interface{} {
	// 获取当前调用栈的信息
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("Failed to retrieve caller information")
	}

	return map[string]interface{}{
		"line": line,
		"file": file,
	}
}
