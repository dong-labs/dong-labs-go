// Package cmd 提供 review 命令
package cmd

import (
	"database/sql"
	"time"

	"github.com/dong-labs/think/internal/core/dates"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/dong-labs/think/internal/think/models"
	"github.com/spf13/cobra"
)

var (
	reviewToday  bool
	reviewWeek   bool
	reviewRandom bool
)

// reviewCmd review 命令
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "回顾想法",
	Run: func(cmd *cobra.Command, args []string) {
		today, _ := cmd.Flags().GetBool("today")
		week, _ := cmd.Flags().GetBool("week")
		random, _ := cmd.Flags().GetBool("random")

		database := db.GetDB()

		query := ""
		params := []interface{}{}
		limit := 1

		// 确定查询逻辑
		if today {
			todayStr := time.Now().Format("2006-01-02")
			query = `
				SELECT id, content, tags, priority, status, context, source_agent, note, created_at, updated_at
				FROM thoughts
				WHERE date(created_at) = ?
				ORDER BY RANDOM()
				LIMIT ?
			`
			params = []interface{}{todayStr, limit}
		} else if week {
			weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
			query = `
				SELECT id, content, tags, priority, status, context, source_agent, note, created_at, updated_at
				FROM thoughts
				WHERE date(created_at) >= ?
				ORDER BY RANDOM()
				LIMIT ?
			`
			params = []interface{}{weekAgo, limit}
		} else if random {
			query = `
				SELECT id, content, tags, priority, status, context, source_agent, note, created_at, updated_at
				FROM thoughts
				ORDER BY RANDOM()
				LIMIT ?
			`
			params = []interface{}{limit}
		} else {
			// 默认返回全部回顾信息（兼容 Go 原有逻辑）
			result := models.ReviewResponse{}

			randomThoughts, _ := database.Query(`
				SELECT id, content, tags, priority, status, context, source_agent, note, created_at, updated_at
				FROM thoughts
				ORDER BY RANDOM()
				LIMIT 5
			`)
			result.Random = scanThoughts(randomThoughts)

			recentThoughts, _ := database.Query(`
				SELECT id, content, tags, priority, status, context, source_agent, note, created_at, updated_at
				FROM thoughts
				ORDER BY created_at DESC
				LIMIT 10
			`)
			result.Recent = scanThoughts(recentThoughts)

			weekStart, weekEnd := dates.ThisWeek(dates.WeekStartMonday)
			weekThoughts, _ := database.Query(`
				SELECT id, content, tags, priority, status, context, source_agent, note, created_at, updated_at
				FROM thoughts
				WHERE date(created_at) >= ? AND date(created_at) <= ?
				ORDER BY created_at DESC
			`, weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))
			result.ThisWeek = scanThoughts(weekThoughts)

			monthStart, monthEnd := dates.ThisMonth()
			monthThoughts, _ := database.Query(`
				SELECT id, content, tags, priority, status, context, source_agent, note, created_at, updated_at
				FROM thoughts
				WHERE date(created_at) >= ? AND date(created_at) <= ?
				ORDER BY created_at DESC
			`, monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02"))
			result.ThisMonth = scanThoughts(monthThoughts)

			output.PrintJSON(result)
			return
		}

		// 执行单条查询
		rows, err := database.Query(query, params...)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		thoughts := scanThoughts(rows)

		if len(thoughts) == 0 {
			message := "还没有想法记录，快去记一条吧！"
			if today {
				message = "今天还没有想法记录"
			} else if week {
				message = "本周还没有想法记录"
			}
			output.PrintJSON(map[string]interface{}{
				"found": false,
				"message": message,
			})
			return
		}

		output.PrintJSON(map[string]interface{}{
			"found": true,
			"thought": thoughts[0],
		})
	},
}

func init() {
	reviewCmd.Flags().BoolVarP(&reviewToday, "today", "", false, "只显示今天的想法")
	reviewCmd.Flags().BoolVarP(&reviewWeek, "week", "", false, "只显示本周的想法")
	reviewCmd.Flags().BoolVarP(&reviewRandom, "random", "", false, "随机显示一条想法")
}

func scanThoughts(rows *sql.Rows) []models.Thought {
	if rows == nil {
		return []models.Thought{}
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
	return thoughts
}
