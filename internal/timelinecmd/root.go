package timelinecmd

import (
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "dong-timeline",
		Short: "线咚咚 - 时间线管理",
		Long:  `所有命令返回 JSON 格式输出，方便 AI 调用。`,
	}
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(AddCmd)
	rootCmd.AddCommand(ListCmd)
	rootCmd.AddCommand(GetCmd)
	rootCmd.AddCommand(DeleteCmd)
	rootCmd.AddCommand(UpdateCmd)
	rootCmd.AddCommand(SearchCmd)
	rootCmd.AddCommand(StatsCmd)
	rootCmd.AddCommand(ExportCmd)
	rootCmd.AddCommand(ImportCmd)
	rootCmd.Execute()
}
