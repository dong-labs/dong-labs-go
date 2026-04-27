package didacmd

import (
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "dong-dida",
		Short: "待咚咚 - 待办事项管理",
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
	rootCmd.AddCommand(CompleteCmd)
	rootCmd.AddCommand(ExportCmd)
	rootCmd.AddCommand(ImportCmd)
	rootCmd.Execute()
}
