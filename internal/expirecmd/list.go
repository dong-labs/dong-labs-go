package expirecmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-expire/db"
	"github.com/dong-labs/think/internal/dong-expire/models"
	"github.com/spf13/cobra"
)

var listLimit int
var listCategory string

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出订阅项",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		category, _ := cmd.Flags().GetString("category")

		database := db.GetDB()
		conn, err := database.GetConnection()
		if err != nil {
			output.PrintJSONError("DB_ERROR", err.Error())
			return
		}

		query := `SELECT id, name, category, expire_date, reminder_days, tags, notes, created_at, updated_at FROM items WHERE 1=1`
		queryArgs := []interface{}{}

		if category != "" {
			query += " AND category = ?"
			queryArgs = append(queryArgs, category)
		}

		query += " ORDER BY expire_date ASC LIMIT ?"
		queryArgs = append(queryArgs, limit)

		rows, err := conn.Query(query, queryArgs...)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		items := []models.Item{}
		for rows.Next() {
			var i models.Item
			var createdAt, updatedAt string
			var reminderDays int
			rows.Scan(&i.ID, &i.Name, &i.Category, &i.ExpireDate, &reminderDays, &i.Tags, &i.Notes, &createdAt, &updatedAt)
			i.ReminderDays = reminderDays
			i.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			i.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
			items = append(items, i)
		}
		output.PrintJSON(map[string]interface{}{"total": len(items), "items": items})
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "返回数量")
	ListCmd.Flags().StringVarP(&listCategory, "category", "c", "", "按分类筛选")
}
