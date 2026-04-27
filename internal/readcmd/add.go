package readcmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-read/db"
	"github.com/spf13/cobra"
)

var addURL, addTitle, addContent, addNote, addSource, addType, addTags string

var AddCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "添加阅读项",
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		url, _ := cmd.Flags().GetString("url")
		_, _ = cmd.Flags().GetString("content")
		note, _ := cmd.Flags().GetString("note")
		source, _ := cmd.Flags().GetString("source")
		typeVal, _ := cmd.Flags().GetString("type")
		tags, _ := cmd.Flags().GetString("tags")

		database := db.GetDB()
		result, err := database.Exec(`INSERT INTO items (url, title, note, source, type_val, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			url, title, note, source, typeVal, tags, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			output.PrintJSONError("INSERT_ERROR", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		output.PrintJSON(map[string]interface{}{"id": id, "title": title})
	},
}

func init() {
	AddCmd.Flags().StringVarP(&addURL, "url", "u", "", "URL链接")
	AddCmd.Flags().StringVarP(&addContent, "content", "c", "", "摘录内容")
	AddCmd.Flags().StringVarP(&addNote, "note", "n", "", "笔记")
	AddCmd.Flags().StringVarP(&addSource, "source", "s", "", "来源")
	AddCmd.Flags().StringVarP(&addType, "type", "t", "quote", "类型 (quote/article/book)")
	AddCmd.Flags().StringVarP(&addTags, "tags", "T", "", "标签")
}
