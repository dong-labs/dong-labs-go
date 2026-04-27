// Package cmd 提供 list 命令
package cmd

import (
	"time"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/core/dates"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/dong-labs/think/internal/think/models"
	"github.com/spf13/cobra"
)

var (
	listLimit     int
	listTag       string
	listPriority  string
	listStatus    string
	listToday     bool
	listWeek      bool
)

// listCmd list 命令
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出想法",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		tag, _ := cmd.Flags().GetString("tag")
		priority, _ := cmd.Flags().GetString("priority")
		status, _ := cmd.Flags().GetString("status")
		today, _ := cmd.Flags().GetBool("today")
		week, _ := cmd.Flags().GetBool("week")

		database := db.GetDB()

		// 构建查询
		query := "SELECT id, content, tags, priority, status, context, source_agent, note, created_at, updated_at FROM thoughts WHERE 1=1"
		var params []interface{}

		if tag != "" {
			query += " AND tags LIKE ?"
			params = append(params, "%"+tag+"%")
		}
		if priority != "" {
			query += " AND priority = ?"
			params = append(params, priority)
		}
		if status != "" {
			query += " AND status = ?"
			params = append(params, status)
		}
		if today {
			todayStr := time.Now().Format("2006-01-02")
			query += " AND date(created_at) = ?"
			params = append(params, todayStr)
		}
		if week {
			weekStart, _ := dates.ThisWeek(dates.WeekStartMonday)
			query += " AND date(created_at) >= ?"
			params = append(params, weekStart.Format("2006-01-02"))
		}

		query += " ORDER BY created_at DESC LIMIT ?"
		params = append(params, limit)

		rows, err := database.Query(query, params...)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		thoughts := []models.Thought{}
		for rows.Next() {
			var t models.Thought
			var createdAt, updatedAt string
			err := rows.Scan(&t.ID, &t.Content, &t.Tags, &t.Priority, &t.Status, &t.Context, &t.SourceAgent, &t.Note, &createdAt, &updatedAt)
			if err != nil {
				continue
			}
			t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
			thoughts = append(thoughts, t)
		}

		output.PrintJSON(map[string]interface{}{
			"total": len(thoughts),
			"items": thoughts,
		})
	},
}

func init() {
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "返回数量")
	listCmd.Flags().StringVarP(&listTag, "tag", "t", "", "按标签筛选")
	listCmd.Flags().StringVarP(&listPriority, "priority", "p", "", "按优先级筛选")
	listCmd.Flags().StringVarP(&listStatus, "status", "s", "", "按状态筛选")
	listCmd.Flags().BoolVarP(&listToday, "today", "", false, "只显示今天的")
	listCmd.Flags().BoolVarP(&listWeek, "week", "", false, "只显示本周的")
}
