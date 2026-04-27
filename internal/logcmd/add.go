package logcmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-log/db"
	"github.com/spf13/cobra"
)

var addContent, addGroup, addDate, addTags string

var AddCmd = &cobra.Command{
	Use:   "add <content>",
	Short: "添加日志",
	Run: func(cmd *cobra.Command, args []string) {
		content := args[0]
		group, _ := cmd.Flags().GetString("group")
		date, _ := cmd.Flags().GetString("date")
		tags, _ := cmd.Flags().GetString("tags")

		// Default values
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}
		if group == "" {
			group = "default"
		}

		database := db.GetDB()
		result, err := database.Exec(`INSERT INTO logs (content, log_group, date, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			content, group, date, tags, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			output.PrintJSONError("INSERT_ERROR", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		output.PrintJSON(map[string]interface{}{"id": id, "content": content, "date": date})
	},
}

func init() {
	AddCmd.Flags().StringVarP(&addGroup, "group", "g", "", "日志组")
	AddCmd.Flags().StringVar(&addDate, "date", "", "日期 YYYY-MM-DD")
	AddCmd.Flags().StringVarP(&addTags, "tags", "t", "", "标签")
}
