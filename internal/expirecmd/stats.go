package expirecmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-expire/db"
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
	today := time.Now().Format("2006-01-02")

	// Total
	var total int
	conn.QueryRow("SELECT COUNT(*) FROM items").Scan(&total)
	stats["total"] = total

	// Expiring soon (within 7 days)
	var expiringSoon int
	conn.QueryRow("SELECT COUNT(*) FROM items WHERE DATE(expire_date) BETWEEN DATE(?) AND DATE(?, '+7 days')", today, today).Scan(&expiringSoon)
	stats["expiring_soon"] = expiringSoon

	// Expired
	var expired int
	conn.QueryRow("SELECT COUNT(*) FROM items WHERE DATE(expire_date) < DATE(?)", today).Scan(&expired)
	stats["expired"] = expired

	// By category
	rows, err := conn.Query(`
		SELECT category, COUNT(*) as count
		FROM items
		GROUP BY category
		ORDER BY count DESC
	`)
	if err == nil {
		defer rows.Close()
		byCategory := []map[string]interface{}{}
		for rows.Next() {
			var category string
			var count int
			rows.Scan(&category, &count)
			byCategory = append(byCategory, map[string]interface{}{"category": category, "count": count})
		}
		stats["by_category"] = byCategory
	}

	return stats, nil
}
