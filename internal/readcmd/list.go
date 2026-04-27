package readcmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-read/db"
	"github.com/dong-labs/think/internal/dong-read/models"
	"github.com/spf13/cobra"
)

var listLimit int
var listType, listSource, listTag string

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出阅读项",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		itemType, _ := cmd.Flags().GetString("type")
		source, _ := cmd.Flags().GetString("source")
		tag, _ := cmd.Flags().GetString("tag")

		database := db.GetDB()
		conn, err := database.GetConnection()
		if err != nil {
			output.PrintJSONError("DB_ERROR", err.Error())
			return
		}

		query := `SELECT id, url, title, note, source, type_val, tags, created_at, updated_at FROM items WHERE 1=1`
		queryArgs := []interface{}{}

		if itemType != "" {
			query += " AND type_val = ?"
			queryArgs = append(queryArgs, itemType)
		}
		if source != "" {
			query += " AND source = ?"
			queryArgs = append(queryArgs, source)
		}
		if tag != "" {
			query += " AND tags LIKE ?"
			queryArgs = append(queryArgs, "%"+tag+"%")
		}

		query += " ORDER BY created_at DESC LIMIT ?"
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
			rows.Scan(&i.ID, &i.URL, &i.Title, &i.Note, &i.Source, &i.Type, &i.Tags, &createdAt, &updatedAt)
			i.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			i.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
			items = append(items, i)
		}
		output.PrintJSON(map[string]interface{}{"total": len(items), "items": items})
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "返回数量")
	ListCmd.Flags().StringVarP(&listType, "type", "t", "", "按类型筛选")
	ListCmd.Flags().StringVarP(&listSource, "source", "s", "", "按来源筛选")
	ListCmd.Flags().StringVarP(&listTag, "tag", "T", "", "按标签筛选")
}
