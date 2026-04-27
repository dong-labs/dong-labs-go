package timelinecmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/timeline/db"
	"github.com/dong-labs/think/internal/timeline/models"
	"github.com/spf13/cobra"
)

var listLimit int

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出事件",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		database := db.GetDB()
		
		conn, err := database.GetConnection()
		if err != nil {
			output.PrintJSONError("DB_ERROR", err.Error())
			return
		}
		
		rows, err := conn.Query(`SELECT id, title, date, description, category, tags, created_at, updated_at FROM events ORDER BY date DESC LIMIT ?`, limit)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		events := []models.Event{}
		for rows.Next() {
			var e models.Event
			var createdAt, updatedAt string
			rows.Scan(&e.ID, &e.Title, &e.Date, &e.Description, &e.Category, &e.Tags, &createdAt, &updatedAt)
			e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			e.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
			events = append(events, e)
		}
		output.PrintJSON(map[string]interface{}{"total": len(events), "items": events})
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "返回数量")
}
