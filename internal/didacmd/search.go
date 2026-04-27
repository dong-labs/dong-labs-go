package didacmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dida/db"
	"github.com/spf13/cobra"
)

var searchStatus, searchPriority string

var SearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "搜索待办",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		status := searchStatus
		priority := searchPriority

		results, err := SearchTodos(keyword, status, priority)
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

func SearchTodos(keyword, status, priority string) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	query := `
		SELECT id, title, content, status, priority, due_date, tags, created_at, updated_at
		FROM todos
		WHERE (title LIKE ? OR content LIKE ?)
	`
	args := []interface{}{"%" + keyword + "%", "%" + keyword + "%"}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if priority != "" {
		query += " AND priority = ?"
		args = append(args, priority)
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
		var title, content, stat, priority, dueDate, tags, createdAt, updatedAt string
		err = rows.Scan(&id, &title, &content, &stat, &priority, &dueDate, &tags, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"id":         id,
			"title":      title,
			"content":    content,
			"status":     stat,
			"priority":   priority,
			"due_date":   dueDate,
			"tags":       parseTags(tags),
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}

	return results, nil
}

func init() {
	SearchCmd.Flags().StringVarP(&searchStatus, "status", "s", "", "按状态搜索")
	SearchCmd.Flags().StringVarP(&searchPriority, "priority", "p", "", "按优先级搜索")
}
