package cangcmd

import (
	"github.com/dong-labs/think/internal/core/output"
	"github.com/spf13/cobra"
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出财务数据",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintJSON(map[string]interface{}{
			"message": "导出功能待实现",
		})
	},
}
