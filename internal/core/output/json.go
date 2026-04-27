// Package output 提供统一的 JSON 输出格式
//
// 所有命令都返回一致格式的 JSON 输出
package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// ErrorInfo 错误信息结构
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Response 统一响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// JsonOutputConfig JSON 输出配置
type JsonOutputConfig struct {
	Indent      bool
	EnsureASCII bool
	SortKeys    bool
}

// DefaultConfig 默认配置
var DefaultConfig = &JsonOutputConfig{
	Indent:      false,
	EnsureASCII: false,
	SortKeys:    false,
}

// PrintJSON 打印 JSON 格式数据
//
// 输出格式: {"success": true, "data": {...}}
func PrintJSON(data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
	}
	printResponse(response, DefaultConfig)
}

// PrintJSONError 打印 JSON 错误格式
//
// 输出格式: {"success": false, "error": {"code": "...", "message": "..."}}
func PrintJSONError(code, message string) {
	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
	printResponse(response, DefaultConfig)
}

// PrintJSONWithError 打印结果或错误
// 如果 err 为 nil，打印 data；否则打印错误信息
func PrintJSONWithError(data interface{}, err error) {
	if err != nil {
		PrintJSONErrorFromError(err)
		return
	}
	PrintJSON(data)
}

// PrintJSONErrorFromError 从 error 对象打印错误
func PrintJSONErrorFromError(err error) {
	// 这里需要导入 errors 包来提取错误信息
	// 为了避免循环依赖，先用简单方式处理
	PrintJSONError("ERROR", err.Error())
}

// PrintJSONWithConfig 使用自定义配置打印 JSON
func PrintJSONWithConfig(data interface{}, config *JsonOutputConfig) {
	response := Response{
		Success: true,
		Data:    data,
	}
	printResponse(response, config)
}

// printResponse 内部打印函数
func printResponse(response Response, config *JsonOutputConfig) {
	var output []byte
	var err error

	if config != nil && config.Indent {
		output, err = json.MarshalIndent(response, "", "  ")
	} else {
		output, err = json.Marshal(response)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
		return
	}

	fmt.Println(string(output))
}

// MustMarshal 序列化为 JSON，panic 如果失败
func MustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// MustMarshalIndent 序列化为带缩进的 JSON，panic 如果失败
func MustMarshalIndent(v interface{}) []byte {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return data
}
