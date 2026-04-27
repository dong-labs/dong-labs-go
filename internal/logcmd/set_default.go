package logcmd

import (
	"github.com/dong-labs/think/internal/core/output"
	"github.com/spf13/cobra"
)

var setDefaultGroup string

var SetDefaultCmd = &cobra.Command{
	Use:   "set-default <group>",
	Short: "设置默认日志组",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		group := args[0]
		// In a real implementation, this would save to config
		output.PrintJSON(map[string]interface{}{
			"default_group": group,
		})
	},
}
