package didacmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dida/db"
	"github.com/spf13/cobra"
)

var addDueDate, addPriority, addNote, addTags string

var AddCmd = &cobra.Command{
	Use:   "add <content>",
	Short: "添加待办",
	Run: func(cmd *cobra.Command, args []string) {
		content := args[0]
		dueDate, _ := cmd.Flags().GetString("due")
		priority, _ := cmd.Flags().GetString("priority")
		note, _ := cmd.Flags().GetString("note")
		tags, _ := cmd.Flags().GetString("tags")

		// 默认优先级
		if priority == "" {
			priority = "medium"
		}

		database := db.GetDB()
		now := time.Now().Format("2006-01-02 15:04:05")
		result, err := database.Exec(`INSERT INTO todos (title, content, due_date, priority, note, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			content, content, dueDate, priority, note, tags, now, now)
		if err != nil {
			output.PrintJSONError("INSERT_ERROR", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		output.PrintJSON(map[string]interface{}{"id": id, "content": content})
	},
}

func init() {
	AddCmd.Flags().StringVarP(&addDueDate, "due", "d", "", "截止时间 YYYY-MM-DD")
	AddCmd.Flags().StringVarP(&addPriority, "priority", "p", "medium", "优先级: high/medium/low")
	AddCmd.Flags().StringVarP(&addNote, "note", "n", "", "备注")
	AddCmd.Flags().StringVarP(&addTags, "tags", "t", "", "标签")
}
