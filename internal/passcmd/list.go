package passcmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/pass/db"
	"github.com/dong-labs/think/internal/pass/models"
	"github.com/spf13/cobra"
)

var (
	listLimit   int
	listCategory string
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出密码",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		category, _ := cmd.Flags().GetString("category")

		database := db.GetDB()
		conn, err := database.GetConnection()
		if err != nil {
			output.PrintJSONError("DB_ERROR", err.Error())
			return
		}

		query := `SELECT id, title, username, url, category, tags, notes, created_at, updated_at FROM passwords WHERE 1=1`
		queryArgs := []interface{}{}

		if category != "" {
			query += " AND category = ?"
			queryArgs = append(queryArgs, category)
		}

		query += " ORDER BY created_at DESC LIMIT ?"
		queryArgs = append(queryArgs, limit)

		rows, err := conn.Query(query, queryArgs...)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		passwords := []models.Password{}
		for rows.Next() {
			var p models.Password
			var createdAt, updatedAt string
			rows.Scan(&p.ID, &p.Title, &p.Username, &p.URL, &p.Category, &p.Tags, &p.Notes, &createdAt, &updatedAt)
			p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
			passwords = append(passwords, p)
		}
		output.PrintJSON(map[string]interface{}{"total": len(passwords), "items": passwords})
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "返回数量")
	ListCmd.Flags().StringVarP(&listCategory, "category", "c", "", "按分类筛选")
}
