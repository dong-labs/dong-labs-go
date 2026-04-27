// Package cmd 提供 get 命令
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

// getCmd get 命令
var getCmd = &cobra.Command{
	Use:   "get <idea_id>",
	Short: "获取想法详情",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ideaID, err := strconv.Atoi(args[0])
		if err != nil {
			output.PrintJSONError("INVALID_ID", "无效的想法 ID")
			return
		}

		database := db.GetDB()

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

		output.PrintJSON(t)
	},
}
