package passcmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/pass/db"
	"github.com/spf13/cobra"
)

var (
	addTitle    string
	addUsername string
	addPassword string
	addURL      string
	addCategory string
	addTags     string
	addNotes    string
)

var AddCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "添加密码",
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		url, _ := cmd.Flags().GetString("url")
		category, _ := cmd.Flags().GetString("category")
		tags, _ := cmd.Flags().GetString("tags")
		notes, _ := cmd.Flags().GetString("notes")

		if password == "" {
			output.PrintJSONError("VALIDATION_ERROR", "密码不能为空")
			return
		}

		database := db.GetDB()
		now := time.Now().Format("2006-01-02 15:04:05")
		result, err := database.Exec(`INSERT INTO passwords (title, username, password, url, category, tags, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			title, username, password, url, category, tags, notes, now, now)
		if err != nil {
			output.PrintJSONError("INSERT_ERROR", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		output.PrintJSON(map[string]interface{}{"id": id, "title": title})
	},
}

func init() {
	AddCmd.Flags().StringVarP(&addUsername, "username", "u", "", "用户名")
	AddCmd.Flags().StringVarP(&addPassword, "password", "p", "", "密码")
	AddCmd.Flags().StringVar(&addURL, "url", "", "URL")
	AddCmd.Flags().StringVarP(&addCategory, "category", "c", "", "分类")
	AddCmd.Flags().StringVarP(&addTags, "tags", "t", "", "标签")
	AddCmd.Flags().StringVarP(&addNotes, "notes", "n", "", "备注")
}
