package expirecmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-expire/db"
	"github.com/spf13/cobra"
)

var updateName, updateCategory, updateExpireDate, updateReminderDays, updateTags, updateNotes string

var UpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "更新订阅项",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		updates := make(map[string]interface{})
		if cmd.Flags().Changed("name") {
			updates["name"] = updateName
		}
		if cmd.Flags().Changed("category") {
			updates["category"] = updateCategory
		}
		if cmd.Flags().Changed("expire_date") {
			updates["expire_date"] = updateExpireDate
		}
		if cmd.Flags().Changed("reminder_days") {
			updates["reminder_days"] = updateReminderDays
		}
		if cmd.Flags().Changed("tags") {
			updates["tags"] = updateTags
		}
		if cmd.Flags().Changed("notes") {
			updates["notes"] = updateNotes
		}
		updates["updated_at"] = time.Now().Format("2006-01-02 15:04:05")

		if len(updates) == 1 {
			output.PrintJSONError("VALIDATION_ERROR", "没有指定要更新的字段")
			return
		}

		result, err := UpdateItem(id, updates)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(result)
	},
}

func UpdateItem(id int, updates map[string]interface{}) (map[string]interface{}, error) {
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

	_, err := database.Exec("UPDATE items SET "+setClause+" WHERE id = ?", args...)
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "更新失败: "+err.Error(), nil)
	}

	return GetItem(id)
}

func init() {
	UpdateCmd.Flags().StringVarP(&updateName, "name", "n", "", "名称")
	UpdateCmd.Flags().StringVarP(&updateCategory, "category", "c", "", "分类")
	UpdateCmd.Flags().StringVarP(&updateExpireDate, "expire_date", "e", "", "到期日期 YYYY-MM-DD")
	UpdateCmd.Flags().StringVarP(&updateReminderDays, "reminder_days", "r", "", "提前提醒天数")
	UpdateCmd.Flags().StringVarP(&updateTags, "tags", "t", "", "标签")
	UpdateCmd.Flags().StringVarP(&updateNotes, "notes", "N", "", "备注")
}
