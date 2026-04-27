// Package cmd 提供公共工具函数
package cmd

// isValidPriority 检查优先级是否有效
func isValidPriority(priority string) bool {
	return priority == "low" || priority == "normal" || priority == "high"
}

// isValidStatus 检查状态是否有效
func isValidStatus(status string) bool {
	return status == "active" || status == "completed" || status == "archived"
}
