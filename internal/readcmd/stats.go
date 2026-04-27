package readcmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-read/db"
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
	conn.QueryRow("SELECT COUNT(*) FROM items").Scan(&total)
	stats["total"] = total

	// By type
	rows, err := conn.Query(`
		SELECT type_val, COUNT(*) as count
		FROM items
		GROUP BY type_val
		ORDER BY count DESC
	`)
	if err == nil {
		defer rows.Close()
		byType := []map[string]interface{}{}
		for rows.Next() {
			var t string
			var count int
			rows.Scan(&t, &count)
			byType = append(byType, map[string]interface{}{"type": t, "count": count})
		}
		stats["by_type"] = byType
	}

	return stats, nil
}
