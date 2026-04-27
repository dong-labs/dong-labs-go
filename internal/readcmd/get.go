package readcmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-read/db"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "获取阅读项详情",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		item, err := GetItem(id)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(item)
	},
}

func GetItem(id int) (map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	var url, title, note, source, typeVal, tags, createdAt, updatedAt string
	err = conn.QueryRow(`
		SELECT url, title, note, source, type_val, tags, created_at, updated_at
		FROM items WHERE id = ?
	`, id).Scan(&url, &title, &note, &source, &typeVal, &tags, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.NewNotFoundError("Item", id)
	}

	return map[string]interface{}{
		"id":         id,
		"url":        url,
		"title":      title,
		"note":       note,
		"source":     source,
		"type":       typeVal,
		"tags":       parseTags(tags),
		"created_at": createdAt,
		"updated_at": updatedAt,
	}, nil
}
