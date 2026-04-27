package expirecmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-expire/db"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "获取订阅项详情",
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

	var name, category, expireDate, reminderDays, tags, notes, createdAt, updatedAt string
	err = conn.QueryRow(`
		SELECT name, category, expire_date, reminder_days, tags, notes, created_at, updated_at
		FROM items WHERE id = ?
	`, id).Scan(&name, &category, &expireDate, &reminderDays, &tags, &notes, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.NewNotFoundError("Item", id)
	}

	return map[string]interface{}{
		"id":            id,
		"name":          name,
		"category":      category,
		"expire_date":   expireDate,
		"reminder_days": reminderDays,
		"tags":          parseTags(tags),
		"notes":         notes,
		"created_at":    createdAt,
		"updated_at":    updatedAt,
	}, nil
}
