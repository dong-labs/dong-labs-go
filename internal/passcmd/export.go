package passcmd

import (
	"github.com/dong-labs/think/internal/core/output"
	"github.com/spf13/cobra"
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出密码数据",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintJSON(map[string]interface{}{
			"message": "导出功能待实现（密码数据不应明文导出）",
		})
	},
}
