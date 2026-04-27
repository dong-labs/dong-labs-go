package logcmd

import (
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "dong-log",
		Short: "记咚咚 - 记录日常日志",
		Long:  `所有命令返回 JSON 格式输出，方便 AI 调用。`,
	}
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(AddCmd)
	rootCmd.AddCommand(ListCmd)
	rootCmd.AddCommand(GetCmd)
	rootCmd.AddCommand(DeleteCmd)
	rootCmd.AddCommand(SearchCmd)
	rootCmd.AddCommand(StatsCmd)
	rootCmd.AddCommand(GroupsCmd)
	rootCmd.AddCommand(SetDefaultCmd)
	rootCmd.AddCommand(ExportCmd)
	rootCmd.AddCommand(ImportCmd)
	rootCmd.Execute()
}
