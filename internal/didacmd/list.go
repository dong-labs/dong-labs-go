package didacmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dida/db"
	"github.com/dong-labs/think/internal/dida/models"
	"github.com/spf13/cobra"
)

var listLimit int
var listStatus, listPriority string

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出待办",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		status, _ := cmd.Flags().GetString("status")
		priority, _ := cmd.Flags().GetString("priority")

		database := db.GetDB()
		conn, err := database.GetConnection()
		if err != nil {
			output.PrintJSONError("DB_ERROR", err.Error())
			return
		}

		query := `SELECT id, title, content, status, priority, due_date, tags, created_at, updated_at FROM todos WHERE 1=1`
		queryArgs := []interface{}{}

		if status != "" {
			query += " AND status = ?"
			queryArgs = append(queryArgs, status)
		}
		if priority != "" {
			query += " AND priority = ?"
			queryArgs = append(queryArgs, priority)
		}

		query += " ORDER BY created_at DESC LIMIT ?"
		queryArgs = append(queryArgs, limit)

		rows, err := conn.Query(query, queryArgs...)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		todos := []models.Todo{}
		for rows.Next() {
			var t models.Todo
			var createdAt, updatedAt string
			rows.Scan(&t.ID, &t.Title, &t.Content, &t.Status, &t.Priority, &t.DueDate, &t.Tags, &createdAt, &updatedAt)
			t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
			todos = append(todos, t)
		}
		output.PrintJSON(map[string]interface{}{"total": len(todos), "items": todos})
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "返回数量")
	ListCmd.Flags().StringVarP(&listStatus, "status", "s", "", "按状态筛选: pending/done/cancelled")
	ListCmd.Flags().StringVarP(&listPriority, "priority", "p", "", "按优先级筛选: high/medium/low")
}
