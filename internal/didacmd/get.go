package didacmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dida/db"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "获取待办详情",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		todo, err := GetTodo(id)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(todo)
	},
}

func GetTodo(id int) (map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	var title, content, status, priority, dueDate, tags, createdAt, updatedAt string
	err = conn.QueryRow(`
		SELECT title, content, status, priority, due_date, tags, created_at, updated_at
		FROM todos WHERE id = ?
	`, id).Scan(&title, &content, &status, &priority, &dueDate, &tags, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.NewNotFoundError("Todo", id)
	}

	return map[string]interface{}{
		"id":         id,
		"title":      title,
		"content":    content,
		"status":     status,
		"priority":   priority,
		"due_date":   dueDate,
		"tags":       parseTags(tags),
		"created_at": createdAt,
		"updated_at": updatedAt,
	}, nil
}
