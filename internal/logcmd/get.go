package logcmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-log/db"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "获取日志详情",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		log, err := GetLog(id)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(log)
	},
}

func GetLog(id int) (map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	var content, logGroup, date, tags, createdAt, updatedAt string
	err = conn.QueryRow(`
		SELECT content, log_group, date, tags, created_at, updated_at
		FROM logs WHERE id = ?
	`, id).Scan(&content, &logGroup, &date, &tags, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.NewNotFoundError("Log", id)
	}

	return map[string]interface{}{
		"id":         id,
		"content":    content,
		"log_group":  logGroup,
		"date":       date,
		"tags":       parseTags(tags),
		"created_at": createdAt,
		"updated_at": updatedAt,
	}, nil
}
