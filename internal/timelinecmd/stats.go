package timelinecmd

import (
	"github.com/dong-labs/think/internal/core/dates"
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/timeline/db"
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

type Stats struct {
	Total      int           `json:"total"`
	ThisWeek   int           `json:"this_week"`
	ThisMonth  int           `json:"this_month"`
	ThisYear   int           `json:"this_year"`
	ByCategory []CategoryStat `json:"by_category"`
}

type CategoryStat struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

func GetStats() (*Stats, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	stats := &Stats{}

	// Total
	conn.QueryRow("SELECT COUNT(*) FROM events").Scan(&stats.Total)

	// This week
	weekStart, _ := dates.ThisWeek(dates.WeekStartMonday)
	conn.QueryRow("SELECT COUNT(*) FROM events WHERE date >= ?", weekStart).Scan(&stats.ThisWeek)

	// This month
	monthStart, _ := dates.ThisMonth()
	conn.QueryRow("SELECT COUNT(*) FROM events WHERE date >= ?", monthStart).Scan(&stats.ThisMonth)

	// This year
	yearStart, _ := dates.ThisYear()
	conn.QueryRow("SELECT COUNT(*) FROM events WHERE date >= ?", yearStart).Scan(&stats.ThisYear)

	// By category
	rows, err := conn.Query(`
		SELECT category, COUNT(*) as count
		FROM events
		GROUP BY category
		ORDER BY count DESC
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cat string
			var count int
			rows.Scan(&cat, &count)
			stats.ByCategory = append(stats.ByCategory, CategoryStat{Category: cat, Count: count})
		}
	}

	return stats, nil
}
