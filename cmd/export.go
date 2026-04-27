// Package cmd 提供 export 命令
package cmd

import (
	"encoding/json"
	"os"
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/dong-labs/think/internal/think/models"
	"github.com/spf13/cobra"
)

var (
	exportOutput string
	exportFormat string
)

// exportCmd export 命令
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出数据",
	Run: func(cmd *cobra.Command, args []string) {
		outputFile, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")

		if format != "json" {
			output.PrintJSONError("UNSUPPORTED_FORMAT", "目前只支持 JSON 格式")
			return
		}

		database := db.GetDB()

		rows, err := database.Query(`
			SELECT id, content, tags, created_at, updated_at
			FROM thoughts
			ORDER BY created_at DESC
		`)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		thoughts := []models.Thought{}
		for rows.Next() {
			var t models.Thought
			var createdAt, updatedAt string
			err := rows.Scan(&t.ID, &t.Content, &t.Tags, &createdAt, &updatedAt)
			if err != nil {
				continue
			}
			t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
			thoughts = append(thoughts, t)
		}

		// 写入文件
		data, _ := json.MarshalIndent(thoughts, "", "  ")
		err = os.WriteFile(outputFile, data, 0644)
		if err != nil {
			output.PrintJSONError("WRITE_ERROR", err.Error())
			return
		}

		output.PrintJSON(map[string]interface{}{
			"file":     outputFile,
			"count":    len(thoughts),
			"format":   format,
			"exported": true,
		})
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "think.json", "输出文件")
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "格式: json")
}
