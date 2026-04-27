package logcmd

import (
	"encoding/json"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-log/db"
	"github.com/spf13/cobra"
)

var exportFile string

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出日志数据",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := FetchAllLogs()
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

func FetchAllLogs() ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	rows, err := conn.Query(`
		SELECT id, content, log_group, date, tags, created_at, updated_at
		FROM logs
		ORDER BY date DESC, created_at DESC
	`)
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
	ExportCmd.Flags().StringVarP(&exportFile, "output", "o", "logs.json", "输出文件")
}
