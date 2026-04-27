package readcmd

import (
	"encoding/json"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-read/db"
	"github.com/spf13/cobra"
)

var exportFile string

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出阅读数据",
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
		SELECT id, url, title, note, source, type_val, tags, created_at, updated_at
		FROM items
		ORDER BY created_at DESC
	`)
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
	ExportCmd.Flags().StringVarP(&exportFile, "output", "o", "read.json", "输出文件")
}
