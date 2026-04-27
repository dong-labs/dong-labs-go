package cangcmd

import (
	"github.com/dong-labs/think/internal/core/output"
	"github.com/spf13/cobra"
)

var importMerge bool

var ImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "导入财务数据",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintJSON(map[string]interface{}{
			"message": "导入功能待实现",
		})
	},
}

func init() {
	ImportCmd.Flags().BoolVarP(&importMerge, "merge", "m", false, "合并模式")
}
