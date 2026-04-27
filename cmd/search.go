// Package cmd 提供 search 命令
package cmd

import (
	"time"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/dong-labs/think/internal/think/models"
	"github.com/spf13/cobra"
)

var (
	searchTag      string
	searchPriority string
	searchStatus   string
)

// searchCmd search 命令
var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "搜索想法",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		tag, _ := cmd.Flags().GetString("tag")
		priority, _ := cmd.Flags().GetString("priority")
		status, _ := cmd.Flags().GetString("status")

		database := db.GetDB()

		// 构建查询条件
		conditions := []string{"content LIKE ?"}
		params := []interface{}{"%" + keyword + "%"}

		if tag != "" {
			conditions = append(conditions, "tags LIKE ?")
			params = append(params, "%"+tag+"%")
		}

		if priority != "" {
			if !isValidPriority(priority) {
				output.PrintJSONError("INVALID_PRIORITY", "无效的优先级值")
				return
			}
			conditions = append(conditions, "priority = ?")
			params = append(params, priority)
		}

		if status != "" {
			if !isValidStatus(status) {
				output.PrintJSONError("INVALID_STATUS", "无效的状态值")
				return
			}
			conditions = append(conditions, "status = ?")
			params = append(params, status)
		}

		whereClause := ""
		for i, c := range conditions {
			if i > 0 {
				whereClause += " AND "
			}
			whereClause += c
		}

		query := `
			SELECT id, content, tags, priority, status, context, source_agent, note, created_at, updated_at
			FROM thoughts
			WHERE ` + whereClause + `
			ORDER BY created_at DESC
		`

		rows, err := database.Query(query, params...)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		thoughts := []models.Thought{}
		for rows.Next() {
			var t models.Thought
			var createdAt, updatedAt string
			err := rows.Scan(&t.ID, &t.Content, &t.Tags, &t.Priority, &t.Status, &t.Context, &t.SourceAgent, &t.Note, &createdAt, &updatedAt)
			if err != nil {
				continue
			}
			t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
			thoughts = append(thoughts, t)
		}

		output.PrintJSON(map[string]interface{}{
			"keyword": keyword,
			"total":   len(thoughts),
			"items":   thoughts,
		})
	},
}

func init() {
	searchCmd.Flags().StringVarP(&searchTag, "tag", "", "", "按标签筛选")
	searchCmd.Flags().StringVarP(&searchPriority, "priority", "", "", "按优先级筛选: low/normal/high")
	searchCmd.Flags().StringVarP(&searchStatus, "status", "", "", "按状态筛选")
}
