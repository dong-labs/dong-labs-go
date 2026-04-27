package passcmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/pass/db"
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

type StatsResponse struct {
	Total      int           `json:"total"`
	ByCategory []CategoryStat `json:"by_category,omitempty"`
}

type CategoryStat struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

func GetStats() (*StatsResponse, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	stats := &StatsResponse{}

	// Total
	conn.QueryRow("SELECT COUNT(*) FROM passwords").Scan(&stats.Total)

	// By category
	rows, err := conn.Query(`SELECT category, COUNT(*) as count FROM passwords GROUP BY category ORDER BY count DESC`)
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
