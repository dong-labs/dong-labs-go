package timelinecmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/timeline/db"
	"github.com/spf13/cobra"
)

var addTitle, addDate, addCategory, addDescription, addTags string

var AddCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "添加事件",
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		date, _ := cmd.Flags().GetString("date")
		category, _ := cmd.Flags().GetString("category")
		description, _ := cmd.Flags().GetString("description")
		tags, _ := cmd.Flags().GetString("tags")

		// Default to today if date not provided
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		database := db.GetDB()
		result, err := database.Exec(`INSERT INTO events (title, date, description, category, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			title, date, description, category, tags, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			output.PrintJSONError("INSERT_ERROR", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		output.PrintJSON(map[string]interface{}{"id": id, "title": title, "date": date})
	},
}

func init() {
	AddCmd.Flags().StringVarP(&addDate, "date", "d", "", "日期 YYYY-MM-DD")
	AddCmd.Flags().StringVarP(&addCategory, "category", "c", "", "分类")
	AddCmd.Flags().StringVar(&addDescription, "description", "", "描述")
	AddCmd.Flags().StringVarP(&addTags, "tags", "t", "", "标签")
}
