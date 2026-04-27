package timelinecmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/timeline/db"
	"github.com/spf13/cobra"
)

var searchCategory string

var SearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "搜索事件",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		category := searchCategory

		results, err := SearchEvents(keyword, category)
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

func SearchEvents(keyword, category string) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	query := `
		SELECT id, title, description, category, tags, date, created_at, updated_at
		FROM events
		WHERE (title LIKE ? OR description LIKE ? OR tags LIKE ?)
	`
	args := []interface{}{"%" + keyword + "%", "%" + keyword + "%", "%" + keyword + "%"}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	query += " ORDER BY date DESC, created_at DESC LIMIT 100"

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var title, description, category, tags, date, createdAt, updatedAt string
		err = rows.Scan(&id, &title, &description, &category, &tags, &date, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": description,
			"category":    category,
			"tags":        parseTags(tags),
			"date":        date,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func init() {
	SearchCmd.Flags().StringVarP(&searchCategory, "category", "c", "", "按分类搜索")
}
