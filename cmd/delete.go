// Package cmd 提供 delete 命令
package cmd

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/dong-labs/think/internal/think/models"
	"github.com/spf13/cobra"
)

var (
	deleteForce bool
)

// deleteCmd delete 命令
var deleteCmd = &cobra.Command{
	Use:   "delete <idea_id>",
	Short: "删除想法",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ideaID, err := strconv.Atoi(args[0])
		if err != nil {
			output.PrintJSONError("INVALID_ID", "无效的想法 ID")
			return
		}

		force, _ := cmd.Flags().GetBool("force")

		database := db.GetDB()

		// 先查询是否存在
		var t models.Thought
		var createdAt, updatedAt string

		err = database.QueryRow(`
			SELECT id, content, tags, priority, status, context, source_agent, note, created_at, updated_at
			FROM thoughts WHERE id = ?
		`, ideaID).Scan(&t.ID, &t.Content, &t.Tags, &t.Priority, &t.Status, &t.Context, &t.SourceAgent, &t.Note, &createdAt, &updatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				output.PrintJSONError("NOT_FOUND", "想法不存在")
			} else {
				output.PrintJSONError("QUERY_ERROR", err.Error())
			}
			return
		}

		t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		// 如果不是强制删除，返回确认信息（JSON 格式）
		if !force {
			output.PrintJSON(map[string]interface{}{
				"cancelled": true,
				"message":   "请使用 --force 参数确认删除",
				"idea":      t,
			})
			return
		}

		// 执行删除
		result, err := database.Exec("DELETE FROM thoughts WHERE id = ?", ideaID)
		if err != nil {
			output.PrintJSONError("DELETE_ERROR", err.Error())
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			output.PrintJSONError("NOT_FOUND", "想法不存在")
			return
		}

		output.PrintJSON(map[string]interface{}{
			"deleted": true,
			"id":      ideaID,
		})
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "强制删除，不提示")
}
