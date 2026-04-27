package expirecmd

import (
	"encoding/json"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-expire/db"
	"github.com/spf13/cobra"
)

var exportFile string

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出订阅数据",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := FetchAllItems()
		if err != nil {
			printError(err)
			return
		}

		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			output.PrintJSONError("MARSHAL_ERROR", err.Error())
			return
		}

		output.PrintJSON(map[string]interface{}{
			"exported": len(data),
			"data":     string(b),
		})
	},
}

func FetchAllItems() ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	rows, err := conn.Query(`
		SELECT id, name, category, expire_date, reminder_days, tags, notes, created_at, updated_at
		FROM items
		ORDER BY expire_date ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name, category, expireDate, reminderDays, tags, notes, createdAt, updatedAt string
		err = rows.Scan(&id, &name, &category, &expireDate, &reminderDays, &tags, &notes, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"id":            id,
			"name":          name,
			"category":      category,
			"expire_date":   expireDate,
			"reminder_days": reminderDays,
			"tags":          parseTags(tags),
			"notes":         notes,
			"created_at":    createdAt,
			"updated_at":    updatedAt,
		})
	}

	return results, nil
}

func init() {
	ExportCmd.Flags().StringVarP(&exportFile, "output", "o", "expire.json", "输出文件")
}
