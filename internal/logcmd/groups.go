package logcmd

import (
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-log/db"
	"github.com/spf13/cobra"
)

var GroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "列出所有日志组",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.GetDB()
		conn, err := database.GetConnection()
		if err != nil {
			output.PrintJSONError("DB_ERROR", err.Error())
			return
		}

		rows, err := conn.Query(`
			SELECT log_group, COUNT(*) as count
			FROM logs
			GROUP BY log_group
			ORDER BY count DESC
		`)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		groups := []map[string]interface{}{}
		for rows.Next() {
			var group string
			var count int
			rows.Scan(&group, &count)
			groups = append(groups, map[string]interface{}{
				"group":  group,
				"count":  count,
			})
		}

		// Get default group from config or use "default"
		defaultGroup := "default"
		
		output.PrintJSON(map[string]interface{}{
			"total":          len(groups),
			"groups":         groups,
			"default_group":  defaultGroup,
		})
	},
}
