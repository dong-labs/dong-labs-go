package passcmd

import (
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/pass/db"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化数据库",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.InitDatabase(); err != nil {
			output.PrintJSONError("INIT_ERROR", err.Error())
			return
		}
		output.PrintJSON(map[string]interface{}{"message": "数据库初始化成功"})
	},
}
