package didacmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dida/db"
	"github.com/spf13/cobra"
)

var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "统计信息",
	Run: func(cmd *cobra.Command, args []string) {
		stats, err := GetStats()
		if err != nil {
			printError(err)
			return
		}
		output.PrintJSON(stats)
	},
}

func GetStats() (map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	stats := make(map[string]interface{})

	// Total
	var total int
	conn.QueryRow("SELECT COUNT(*) FROM todos").Scan(&total)
	stats["total"] = total

	// By status
	rows, err := conn.Query(`
		SELECT status, COUNT(*) as count
		FROM todos
		GROUP BY status
	`)
	if err == nil {
		defer rows.Close()
		byStatus := []map[string]interface{}{}
		for rows.Next() {
			var status string
			var count int
			rows.Scan(&status, &count)
			byStatus = append(byStatus, map[string]interface{}{"status": status, "count": count})
		}
		stats["by_status"] = byStatus
	}

	// By priority
	rows2, err := conn.Query(`
		SELECT priority, COUNT(*) as count
		FROM todos
		GROUP BY priority
	`)
	if err == nil {
		defer rows2.Close()
		byPriority := []map[string]interface{}{}
		for rows2.Next() {
			var priority string
			var count int
			rows2.Scan(&priority, &count)
			byPriority = append(byPriority, map[string]interface{}{"priority": priority, "count": count})
		}
		stats["by_priority"] = byPriority
	}

	return stats, nil
}
