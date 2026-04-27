package timelinecmd

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/timeline/db"
	"github.com/spf13/cobra"
)

var (
	importMerge bool
)

var ImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "导入事件数据",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		content, err := os.ReadFile(file)
		if err != nil {
			output.PrintJSONError("READ_ERROR", err.Error())
			return
		}

		var data []map[string]interface{}

		// Detect format by extension
		if strings.HasSuffix(file, ".json") {
			data, err = ImportFromJSON(string(content))
		} else if strings.HasSuffix(file, ".csv") {
			data, err = ImportFromCSV(string(content))
		} else {
			output.PrintJSONError("VALIDATION_ERROR", "不支持的文件格式")
			return
		}

		if err != nil {
			printError(err)
			return
		}

		if !importMerge {
			// Clear existing data
			db.GetDB().Exec("DELETE FROM events")
		}

		added := 0
		for _, item := range data {
			_, err := AddEventFromImport(item)
			if err != nil {
				continue
			}
			added++
		}

		output.PrintJSON(map[string]interface{}{
			"added": added,
		})
	},
}

func AddEventFromImport(data map[string]interface{}) (int, error) {
	database := db.GetDB()

	title := ""
	description := ""
	category := ""
	tags := []string{}

	if v, ok := data["title"].(string); ok {
		title = v
	}
	if v, ok := data["description"].(string); ok {
		description = v
	}
	if v, ok := data["category"].(string); ok {
		category = v
	}
	if v, ok := data["tags"].([]interface{}); ok {
		for _, tag := range v {
			if s, ok := tag.(string); ok {
				tags = append(tags, s)
			}
		}
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := database.Exec(
		`INSERT INTO events (title, description, category, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		title, description, category, strings.Join(tags, ","), now, now)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func ImportFromJSON(jsonStr string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	return data, err
}

func ImportFromCSV(csvStr string) ([]map[string]interface{}, error) {
	r := csv.NewReader(strings.NewReader(csvStr))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return []map[string]interface{}{}, nil
	}

	// Parse header
	headers := records[0]
	results := []map[string]interface{}{}

	for i := 1; i < len(records); i++ {
		row := records[i]
		item := make(map[string]interface{})
		for j, val := range row {
			if j < len(headers) {
				item[headers[j]] = val
			}
		}
		results = append(results, item)
	}

	return results, nil
}

func init() {
	ImportCmd.Flags().BoolVarP(&importMerge, "merge", "m", false, "合并模式（追加而非替换）")
}
