// Package cmd 提供 stats 命令
package cmd

import (
	"github.com/dong-labs/think/internal/core/dates"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/dong-labs/think/internal/think/models"
	"github.com/spf13/cobra"
)

// statsCmd stats 命令
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "统计信息",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.GetDB()

		result := models.StatsResponse{}

		// 总数
		database.QueryRow("SELECT COUNT(*) FROM thoughts").Scan(&result.Total)

		// 本周
		weekStart, weekEnd := dates.ThisWeek(dates.WeekStartMonday)
		database.QueryRow(`
			SELECT COUNT(*) FROM thoughts
			WHERE date(created_at) >= ? AND date(created_at) <= ?
		`, weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02")).Scan(&result.ThisWeek)

		// 本月
		monthStart, monthEnd := dates.ThisMonth()
		database.QueryRow(`
			SELECT COUNT(*) FROM thoughts
			WHERE date(created_at) >= ? AND date(created_at) <= ?
		`, monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02")).Scan(&result.ThisMonth)

		// 本年
		yearStart, yearEnd := dates.ThisYear()
		database.QueryRow(`
			SELECT COUNT(*) FROM thoughts
			WHERE date(created_at) >= ? AND date(created_at) <= ?
		`, yearStart.Format("2006-01-02"), yearEnd.Format("2006-01-02")).Scan(&result.ThisYear)

		// 按状态统计
		statusStats := make(map[string]int)
		rows, _ := database.Query(`
			SELECT status, COUNT(*) as count
			FROM thoughts
			GROUP BY status
		`)
		defer rows.Close()
		for rows.Next() {
			var status string
			var count int
			rows.Scan(&status, &count)
			statusStats[status] = count
		}
		result.StatusStats = statusStats

		// 按优先级统计
		priorityStats := make(map[string]int)
		rows2, _ := database.Query(`
			SELECT priority, COUNT(*) as count
			FROM thoughts
			GROUP BY priority
		`)
		defer rows2.Close()
		for rows2.Next() {
			var priority string
			var count int
			rows2.Scan(&priority, &count)
			priorityStats[priority] = count
		}
		result.PriorityStats = priorityStats

		// 标签统计
		rows3, _ := database.Query(`
			SELECT tags, COUNT(*) as count
			FROM thoughts
			WHERE tags IS NOT NULL AND tags != ''
			GROUP BY tags
			ORDER BY count DESC
			LIMIT 10
		`)
		defer rows3.Close()

		topTags := []models.TagStat{}
		for rows3.Next() {
			var t models.TagStat
			rows3.Scan(&t.Tag, &t.Count)
			topTags = append(topTags, t)
		}
		result.TopTags = topTags

		output.PrintJSON(result)
	},
}
