package timelinecmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/timeline/db"
	"github.com/spf13/cobra"
)

var (
	exportFile   string
	exportFormat string
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出事件数据",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := FetchAllEvents()
		if err != nil {
			printError(err)
			return
		}

		var content string
		switch exportFormat {
		case "json":
			content, err = ExportToJSON(data)
		case "csv":
			content, err = ExportToCSV(data)
		case "md", "markdown":
			content, err = ExportToMarkdown(data)
		default:
			output.PrintJSONError("VALIDATION_ERROR", "不支持的格式: "+exportFormat)
			return
		}

		if err != nil {
			printError(err)
			return
		}

		if err := os.WriteFile(exportFile, []byte(content), 0644); err != nil {
			output.PrintJSONError("WRITE_ERROR", err.Error())
			return
		}

		output.PrintJSON(map[string]interface{}{
			"exported": len(data),
			"file":     exportFile,
			"format":   exportFormat,
		})
	},
}

func FetchAllEvents() ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	rows, err := conn.Query(`
		SELECT id, title, description, category, tags, date, created_at, updated_at
		FROM events
		ORDER BY date DESC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var title, description, category, tags, date, createdAt, updatedAt string
		err = rows.Scan(&id, &title, &description, &category, &tags, &date, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": description,
			"category":    category,
			"tags":        parseTags(tags),
			"date":        date,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
	}

	return results, nil
}

func ExportToJSON(data []map[string]interface{}) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ExportToCSV(data []map[string]interface{}) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	var buf strings.Builder
	w := csv.NewWriter(&buf)

	// Header
	headers := []string{"id", "title", "description", "category", "tags", "date", "created_at", "updated_at"}
	w.Write(headers)

	// Rows
	for _, item := range data {
		tags := ""
		if t, ok := item["tags"].([]string); ok {
			tags = strings.Join(t, ",")
		}
		w.Write([]string{
			fmt.Sprintf("%v", item["id"]),
			fmt.Sprintf("%v", item["title"]),
			fmt.Sprintf("%v", item["description"]),
			fmt.Sprintf("%v", item["category"]),
			tags,
			fmt.Sprintf("%v", item["date"]),
			fmt.Sprintf("%v", item["created_at"]),
			fmt.Sprintf("%v", item["updated_at"]),
		})
	}

	w.Flush()
	return buf.String(), nil
}

func ExportToMarkdown(data []map[string]interface{}) (string, error) {
	var buf strings.Builder
	buf.WriteString("# Timeline 数据导出\n\n")

	for _, item := range data {
		title := item["title"]
		date := item["date"]
		description := item["description"]
		category := item["category"]

		buf.WriteString(fmt.Sprintf("## %s (%s)\n\n", title, date))
		if category != "" {
			buf.WriteString(fmt.Sprintf("**分类**: %s\n\n", category))
		}
		if description != "" {
			buf.WriteString(fmt.Sprintf("%s\n\n", description))
		}
		buf.WriteString("---\n\n")
	}

	return buf.String(), nil
}

func init() {
	ExportCmd.Flags().StringVarP(&exportFile, "output", "o", "timeline.json", "输出文件")
	ExportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "格式: json/csv/md")
}
