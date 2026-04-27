package cangcmd

import (
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "dong-cang",
		Short: "财咚咚 - 财务管理",
		Long:  `所有命令返回 JSON 格式输出，方便 AI 调用。`,
	}
	
	// 主要命令
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(TxCmd)
	rootCmd.AddCommand(ListCmd)
	rootCmd.AddCommand(GetCmd)
	rootCmd.AddCommand(DeleteCmd)
	rootCmd.AddCommand(SearchCmd)
	rootCmd.AddCommand(StatsCmd)
	rootCmd.AddCommand(ExportCmd)
	rootCmd.AddCommand(ImportCmd)
	
	// account 子命令
	accountCmd := &cobra.Command{
		Use:   "account",
		Short: "账户管理",
	}
	accountCmd.AddCommand(AccountAddCmd)
	rootCmd.AddCommand(accountCmd)
	
	rootCmd.Execute()
}
