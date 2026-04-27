package passcmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/pass/db"
	"github.com/spf13/cobra"
)

var searchCategory string

var SearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "搜索密码",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		category := searchCategory

		results, err := SearchPasswords(keyword, category)
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

func SearchPasswords(keyword, category string) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	query := `SELECT id, title, username, url, category, tags, notes, created_at, updated_at FROM passwords WHERE (title LIKE ? OR username LIKE ? OR url LIKE ? OR notes LIKE ?)`
	args := []interface{}{"%" + keyword + "%", "%" + keyword + "%", "%" + keyword + "%", "%" + keyword + "%"}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var title, username, url, cat, tags, notes, createdAt, updatedAt string
		err = rows.Scan(&id, &title, &username, &url, &cat, &tags, &notes, &createdAt, &updatedAt)
		if err != nil { continue }

		results = append(results, map[string]interface{}{
			"id":         id,
			"title":      title,
			"username":   username,
			"url":        url,
			"category":   cat,
			"tags":       tags,
			"notes":      notes,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}

	return results, nil
}

func init() {
	SearchCmd.Flags().StringVarP(&searchCategory, "category", "c", "", "按分类搜索")
}
