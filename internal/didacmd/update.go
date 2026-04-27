package didacmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dida/db"
	"github.com/spf13/cobra"
)

var updateContent, updateStatus, updatePriority, updateDueDate, updateNote, updateTags string

var UpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "更新待办",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		updates := make(map[string]interface{})
		if cmd.Flags().Changed("content") {
			updates["title"] = updateContent
			updates["content"] = updateContent
		}
		if cmd.Flags().Changed("status") {
			updates["status"] = updateStatus
			if updateStatus == "done" {
				updates["completed_at"] = time.Now().Format("2006-01-02 15:04:05")
			}
		}
		if cmd.Flags().Changed("priority") {
			updates["priority"] = updatePriority
		}
		if cmd.Flags().Changed("due") {
			updates["due_date"] = updateDueDate
		}
		if cmd.Flags().Changed("note") {
			updates["note"] = updateNote
		}
		if cmd.Flags().Changed("tags") {
			updates["tags"] = updateTags
		}
		updates["updated_at"] = time.Now().Format("2006-01-02 15:04:05")

		if len(updates) == 1 {
			output.PrintJSONError("VALIDATION_ERROR", "没有指定要更新的字段")
			return
		}

		result, err := UpdateTodo(id, updates)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(result)
	},
}

func UpdateTodo(id int, updates map[string]interface{}) (map[string]interface{}, error) {
	database := db.GetDB()

	setClause := ""
	args := make([]interface{}, 0)
	for key, val := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += key + " = ?"
		args = append(args, val)
	}
	args = append(args, id)

	_, err := database.Exec("UPDATE todos SET "+setClause+" WHERE id = ?", args...)
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "更新失败: "+err.Error(), nil)
	}

	return GetTodo(id)
}

func init() {
	UpdateCmd.Flags().StringVarP(&updateContent, "content", "c", "", "待办内容")
	UpdateCmd.Flags().StringVarP(&updateStatus, "status", "s", "", "状态: pending/done/cancelled")
	UpdateCmd.Flags().StringVarP(&updatePriority, "priority", "p", "", "优先级: high/medium/low")
	UpdateCmd.Flags().StringVarP(&updateDueDate, "due", "d", "", "截止时间 YYYY-MM-DD")
	UpdateCmd.Flags().StringVarP(&updateNote, "note", "n", "", "备注")
	UpdateCmd.Flags().StringVarP(&updateTags, "tags", "t", "", "标签")
}
