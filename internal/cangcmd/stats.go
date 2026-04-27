package cangcmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/cang/db"
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

	// Total transactions
	var total int
	conn.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&total)
	stats["total_transactions"] = total

	// Total amount
	var totalAmount int
	conn.QueryRow("SELECT COALESCE(SUM(amount_cents), 0) FROM transactions").Scan(&totalAmount)
	stats["total_amount"] = float64(totalAmount) / 100

	// This month
	monthStart := today[:8] + "01"
	var thisMonth int
	conn.QueryRow("SELECT COALESCE(SUM(amount_cents), 0) FROM transactions WHERE date >= ?", monthStart).Scan(&thisMonth)
	stats["this_month_amount"] = float64(thisMonth) / 100

	// By category
	rows, err := conn.Query(`SELECT category, COUNT(*) as count, COALESCE(SUM(amount_cents), 0) as amount FROM transactions GROUP BY category ORDER BY amount DESC`)
	if err == nil {
		defer rows.Close()
		byCategory := []map[string]interface{}{}
		for rows.Next() {
			var cat string
			var count, amount int
			rows.Scan(&cat, &count, &amount)
			byCategory = append(byCategory, map[string]interface{}{
				"category": cat,
				"count":    count,
				"amount":   float64(amount) / 100,
			})
		}
		stats["by_category"] = byCategory
	}

	return stats, nil
}
