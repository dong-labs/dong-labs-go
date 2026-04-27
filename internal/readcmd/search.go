package readcmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-read/db"
	"github.com/spf13/cobra"
)

var searchType, searchTag string

var SearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "搜索阅读项",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		itemType := searchType
		tag := searchTag

		results, err := SearchItems(keyword, itemType, tag)
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

func SearchItems(keyword, itemType, tag string) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	query := `
		SELECT id, url, title, note, source, type_val, tags, created_at, updated_at
		FROM items
		WHERE (title LIKE ? OR note LIKE ?)
	`
	args := []interface{}{"%" + keyword + "%", "%" + keyword + "%"}

	if itemType != "" {
		query += " AND type_val = ?"
		args = append(args, itemType)
	}
	if tag != "" {
		query += " AND tags LIKE ?"
		args = append(args, "%"+tag+"%")
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
		var url, title, note, source, typeVal, tags, createdAt, updatedAt string
		err = rows.Scan(&id, &url, &title, &note, &source, &typeVal, &tags, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"id":         id,
			"url":        url,
			"title":      title,
			"note":       note,
			"source":     source,
			"type":       typeVal,
			"tags":       parseTags(tags),
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}

	return results, nil
}

func init() {
	SearchCmd.Flags().StringVarP(&searchType, "type", "t", "", "按类型搜索")
	SearchCmd.Flags().StringVarP(&searchTag, "tag", "T", "", "按标签搜索")
}
