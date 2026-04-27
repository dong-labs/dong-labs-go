// Package cmd 提供 add 命令
package cmd

import (
	"time"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/spf13/cobra"
)

var (
	addTag      string
	addPriority string
	addContext  string
	addSource   string
	addNote     string
)

// addCmd add 命令
var addCmd = &cobra.Command{
	Use:   "add <content>",
	Short: "记录想法",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content := args[0]
		tag, _ := cmd.Flags().GetString("tag")
		priority, _ := cmd.Flags().GetString("priority")
		context, _ := cmd.Flags().GetString("context")
		source, _ := cmd.Flags().GetString("source")
		note, _ := cmd.Flags().GetString("note")

		database := db.GetDB()

		// 插入数据
		result, err := database.Exec(`
			INSERT INTO thoughts (content, tags, priority, context, source_agent, note, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, content, tag, priority, context, source, note, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

		if err != nil {
			output.PrintJSONError("INSERT_ERROR", err.Error())
			return
		}

		id, _ := result.LastInsertId()

		output.PrintJSON(map[string]interface{}{
			"id":          id,
			"content":     content,
			"tags":        tag,
			"priority":    priority,
			"context":     context,
			"source_agent": source,
			"note":        note,
		})
	},
}

func init() {
	addCmd.Flags().StringVarP(&addTag, "tag", "t", "", "标签")
	addCmd.Flags().StringVarP(&addPriority, "priority", "p", "normal", "优先级: low/normal/high")
	addCmd.Flags().StringVarP(&addContext, "context", "c", "", "上下文")
	addCmd.Flags().StringVarP(&addSource, "source", "s", "", "来源智能体")
	addCmd.Flags().StringVarP(&addNote, "note", "n", "", "备注")
}
