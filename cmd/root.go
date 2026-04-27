// Package cmd 提供 think-cli 的命令入口
package cmd

import (
	"github.com/spf13/cobra"
)

// Execute 执行 think 命令
func Execute() {
	rootCmd := &cobra.Command{
		Use:   "dong-think",
		Short: "思咚咚 - 记录灵感和想法",
		Long:  `所有命令返回 JSON 格式输出，方便 AI 调用。`,
	}
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(tagsCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.Execute()
}
