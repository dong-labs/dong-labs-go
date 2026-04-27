package logcmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-log/db"
	"github.com/spf13/cobra"
)

var searchGroup, searchTag string

var SearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "搜索日志",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		group := searchGroup
		tag := searchTag

		results, err := SearchLogs(keyword, group, tag)
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

func SearchLogs(keyword, group, tag string) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	query := `
		SELECT id, content, log_group, date, tags, created_at, updated_at
		FROM logs
		WHERE (content LIKE ? OR tags LIKE ?)
	`
	args := []interface{}{"%" + keyword + "%", "%" + keyword + "%"}

	if group != "" {
		query += " AND log_group = ?"
		args = append(args, group)
	}
	if tag != "" {
		query += " AND tags LIKE ?"
		args = append(args, "%"+tag+"%")
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
		var content, logGroup, date, tags, createdAt, updatedAt string
		err = rows.Scan(&id, &content, &logGroup, &date, &tags, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"id":         id,
			"content":    content,
			"log_group":  logGroup,
			"date":       date,
			"tags":       parseTags(tags),
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}

	return results, nil
}

func init() {
	SearchCmd.Flags().StringVarP(&searchGroup, "group", "g", "", "按组搜索")
	SearchCmd.Flags().StringVarP(&searchTag, "tag", "t", "", "按标签搜索")
}
