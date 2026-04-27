package expirecmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-expire/db"
	"github.com/spf13/cobra"
)

var searchCategory string

var SearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "搜索订阅项",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		category := searchCategory

		results, err := SearchItems(keyword, category)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(map[string]interface{}{
			"total": len(results),
			"items": results,
		})
	},
}

func SearchItems(keyword, category string) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	query := `
		SELECT id, name, category, expire_date, reminder_days, tags, notes, created_at, updated_at
		FROM items
		WHERE (name LIKE ? OR notes LIKE ?)
	`
	args := []interface{}{"%" + keyword + "%", "%" + keyword + "%"}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	query += " ORDER BY expire_date ASC LIMIT 100"

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name, category, expireDate, reminderDays, tags, notes, createdAt, updatedAt string
		err = rows.Scan(&id, &name, &category, &expireDate, &reminderDays, &tags, &notes, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"id":            id,
			"name":          name,
			"category":      category,
			"expire_date":   expireDate,
			"reminder_days": reminderDays,
			"tags":          parseTags(tags),
			"notes":         notes,
			"created_at":    createdAt,
			"updated_at":    updatedAt,
		})
	}

	return results, nil
}

func init() {
	SearchCmd.Flags().StringVarP(&searchCategory, "category", "c", "", "按分类搜索")
}
