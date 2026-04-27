// Package cmd 提供 init 命令
package cmd

import (
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/spf13/cobra"
)

// initCmd init 命令
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化数据库",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.InitDatabase(); err != nil {
			output.PrintJSONError("INIT_ERROR", err.Error())
			return
		}
		output.PrintJSON(map[string]interface{}{
			"message": "数据库初始化成功",
			"version": db.SCHEMA_VERSION,
		})
	},
}
