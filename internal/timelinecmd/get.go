package timelinecmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/timeline/db"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "获取事件详情",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		event, err := GetEvent(id)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(event)
	},
}

func GetEvent(id int) (map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	var title, description, category, tags, date, createdAt, updatedAt string
	err = conn.QueryRow(`
		SELECT title, description, category, tags, date, created_at, updated_at
		FROM events WHERE id = ?
	`, id).Scan(&title, &description, &category, &tags, &date, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.NewNotFoundError("Event", id)
	}

	return map[string]interface{}{
		"id":          id,
		"title":       title,
		"description": description,
		"category":    category,
		"tags":        parseTags(tags),
		"date":        date,
		"created_at":  createdAt,
		"updated_at":  updatedAt,
	}, nil
}
