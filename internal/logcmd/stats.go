package logcmd

import (
	"github.com/dong-labs/think/internal/core/dates"
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-log/db"
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
	conn.QueryRow("SELECT COUNT(*) FROM logs").Scan(&total)
	stats["total"] = total

	// This week
	weekStart, _ := dates.ThisWeek(dates.WeekStartMonday)
	var thisWeek int
	conn.QueryRow("SELECT COUNT(*) FROM logs WHERE date >= ?", weekStart).Scan(&thisWeek)
	stats["this_week"] = thisWeek

	// This month
	monthStart, _ := dates.ThisMonth()
	var thisMonth int
	conn.QueryRow("SELECT COUNT(*) FROM logs WHERE date >= ?", monthStart).Scan(&thisMonth)
	stats["this_month"] = thisMonth

	// This year
	yearStart, _ := dates.ThisYear()
	var thisYear int
	conn.QueryRow("SELECT COUNT(*) FROM logs WHERE date >= ?", yearStart).Scan(&thisYear)
	stats["this_year"] = thisYear

	return stats, nil
}
