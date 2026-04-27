package expirecmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-expire/db"
	"github.com/spf13/cobra"
)

var addCategory, addExpireDate, addReminderDays, addTags, addNotes string

var AddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "添加订阅项",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		category, _ := cmd.Flags().GetString("category")
		expireDate, _ := cmd.Flags().GetString("expire_date")
		reminderDays, _ := cmd.Flags().GetString("reminder_days")
		tags, _ := cmd.Flags().GetString("tags")
		notes, _ := cmd.Flags().GetString("notes")

		database := db.GetDB()
		now := time.Now().Format("2006-01-02 15:04:05")
		result, err := database.Exec(`INSERT INTO items (name, category, expire_date, reminder_days, tags, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			name, category, expireDate, reminderDays, tags, notes, now, now)
		if err != nil {
			output.PrintJSONError("INSERT_ERROR", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		output.PrintJSON(map[string]interface{}{"id": id, "name": name, "expire_date": expireDate})
	},
}

func init() {
	AddCmd.Flags().StringVarP(&addCategory, "category", "c", "", "分类")
	AddCmd.Flags().StringVarP(&addExpireDate, "expire_date", "e", "", "到期日期 YYYY-MM-DD")
	AddCmd.Flags().StringVarP(&addReminderDays, "reminder_days", "r", "7", "提前提醒天数")
	AddCmd.Flags().StringVarP(&addTags, "tags", "t", "", "标签")
	AddCmd.Flags().StringVarP(&addNotes, "notes", "n", "", "备注")
}
