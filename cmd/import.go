// Package cmd 提供 import 命令
package cmd

import (
	"encoding/json"
	"os"
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/spf13/cobra"
)

var (
	importFile  string
	importMerge bool
	importDryRun bool
)

// importCmd import 命令
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "导入数据",
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		merge, _ := cmd.Flags().GetBool("merge")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// 读取文件
		data, err := os.ReadFile(file)
		if err != nil {
			output.PrintJSONError("READ_ERROR", err.Error())
			return
		}

		var thoughts []map[string]interface{}
		if err := json.Unmarshal(data, &thoughts); err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		database := db.GetDB()

		added := 0
		skipped := 0
		failed := 0

		for _, t := range thoughts {
			content, _ := t["content"].(string)
			tags, _ := t["tags"].(string)

			if content == "" {
				failed++
				continue
			}

			// 检查是否已存在
			if !merge {
				var exists bool
				database.QueryRow("SELECT EXISTS(SELECT 1 FROM thoughts WHERE content = ?)", content).Scan(&exists)
				if exists {
					skipped++
					continue
				}
			}

			if dryRun {
				added++
				continue
			}

			// 插入
			_, err := database.Exec(`
				INSERT INTO thoughts (content, tags, created_at, updated_at)
				VALUES (?, ?, ?, ?)
			`, content, tags, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

			if err != nil {
				failed++
			} else {
				added++
			}
		}

		output.PrintJSON(map[string]interface{}{
			"total":   len(thoughts),
			"added":   added,
			"skipped": skipped,
			"failed":  failed,
			"dry_run": dryRun,
		})
	},
}

func init() {
	importCmd.Flags().StringVarP(&importFile, "file", "f", "", "导入文件")
	importCmd.Flags().BoolVar(&importMerge, "merge", false, "合并模式")
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "预览模式")

	importCmd.MarkFlagRequired("file")
}
