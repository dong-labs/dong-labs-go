package logcmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-log/db"
	"github.com/dong-labs/think/internal/dong-log/models"
	"github.com/spf13/cobra"
)

var listLimit int
var listGroup, listTag string

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出日志",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		group, _ := cmd.Flags().GetString("group")
		tag, _ := cmd.Flags().GetString("tag")

		database := db.GetDB()
		conn, err := database.GetConnection()
		if err != nil {
			output.PrintJSONError("DB_ERROR", err.Error())
			return
		}

		query := `SELECT id, content, log_group, date, tags, created_at, updated_at FROM logs WHERE 1=1`
		queryArgs := []interface{}{}

		if group != "" {
			query += " AND log_group = ?"
			queryArgs = append(queryArgs, group)
		}
		if tag != "" {
			query += " AND tags LIKE ?"
			queryArgs = append(queryArgs, "%"+tag+"%")
		}

		query += " ORDER BY date DESC, created_at DESC LIMIT ?"
		queryArgs = append(queryArgs, limit)

		rows, err := conn.Query(query, queryArgs...)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		logs := []models.Log{}
		for rows.Next() {
			var l models.Log
			var createdAt, updatedAt string
			rows.Scan(&l.ID, &l.Content, &l.LogGroup, &l.Date, &l.Tags, &createdAt, &updatedAt)
			l.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			l.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
			logs = append(logs, l)
		}
		output.PrintJSON(map[string]interface{}{"total": len(logs), "items": logs})
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "返回数量")
	ListCmd.Flags().StringVarP(&listGroup, "group", "g", "", "按组筛选")
	ListCmd.Flags().StringVarP(&listTag, "tag", "t", "", "按标签筛选")
}
