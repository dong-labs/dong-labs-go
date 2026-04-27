// Package cmd 提供 update 命令
package cmd

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/spf13/cobra"
)

var (
	updateStatus    string
	updatePriority  string
	updateAddTag    string
	updateRemoveTag string
	updateNote      string
)

// updateCmd update 命令
var updateCmd = &cobra.Command{
	Use:   "update <idea_id>",
	Short: "更新想法",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ideaID, err := strconv.Atoi(args[0])
		if err != nil {
			output.PrintJSONError("INVALID_ID", "无效的想法 ID")
			return
		}

		status, _ := cmd.Flags().GetString("status")
		priority, _ := cmd.Flags().GetString("priority")
		tagAdd, _ := cmd.Flags().GetString("add-tag")
		tagRemove, _ := cmd.Flags().GetString("remove-tag")
		note, _ := cmd.Flags().GetString("note")

		database := db.GetDB()

		// 获取当前想法
		var content, tagsStr, currentStatus, currentPriority, currentNote string
		err = database.QueryRow(`
			SELECT content, tags, status, priority, note
			FROM thoughts WHERE id = ?
		`, ideaID).Scan(&content, &tagsStr, &currentStatus, &currentPriority, &currentNote)

		if err != nil {
			if err == sql.ErrNoRows {
				output.PrintJSONError("NOT_FOUND", "想法不存在")
			} else {
				output.PrintJSONError("QUERY_ERROR", err.Error())
			}
			return
		}

		// 解析当前标签
		var currentTags []string
		if tagsStr != "" {
			json.Unmarshal([]byte(tagsStr), &currentTags)
		}
		if currentTags == nil {
			currentTags = []string{}
		}

		// 更新状态
		if status != "" {
			if !isValidStatus(status) {
				output.PrintJSONError("INVALID_STATUS", "无效的状态值")
				return
			}
			currentStatus = status
		}

		// 更新优先级
		if priority != "" {
			if !isValidPriority(priority) {
				output.PrintJSONError("INVALID_PRIORITY", "无效的优先级值")
				return
			}
			currentPriority = priority
		}

		// 添加标签
		if tagAdd != "" {
			newTag := strings.TrimSpace(tagAdd)
			if newTag != "" && !containsTag(currentTags, newTag) {
				currentTags = append(currentTags, newTag)
			}
		}

		// 移除标签
		if tagRemove != "" {
			tagToDelete := strings.TrimSpace(tagRemove)
			currentTags = removeTag(currentTags, tagToDelete)
		}

		// 更新备注
		if note != "" {
			currentNote = note
		}

		// 序列化标签
		newTagsStr := ""
		if len(currentTags) > 0 {
			tagBytes, _ := json.Marshal(currentTags)
			newTagsStr = string(tagBytes)
		}

		// 更新数据库
		_, err = database.Exec(`
			UPDATE thoughts
			SET tags = ?, priority = ?, status = ?, note = ?, updated_at = ?
			WHERE id = ?
		`, newTagsStr, currentPriority, currentStatus, currentNote, time.Now().Format("2006-01-02 15:04:05"), ideaID)

		if err != nil {
			output.PrintJSONError("UPDATE_ERROR", err.Error())
			return
		}

		output.PrintJSON(map[string]interface{}{
			"updated":  ideaID,
			"status":   currentStatus,
			"priority": currentPriority,
			"tags":     currentTags,
			"note":     currentNote,
		})
	},
}

func init() {
	updateCmd.Flags().StringVarP(&updateStatus, "status", "", "", "更新状态")
	updateCmd.Flags().StringVarP(&updatePriority, "priority", "", "", "更新优先级: low/normal/high")
	updateCmd.Flags().StringVarP(&updateAddTag, "add-tag", "", "", "添加标签")
	updateCmd.Flags().StringVarP(&updateRemoveTag, "remove-tag", "", "", "移除标签")
	updateCmd.Flags().StringVarP(&updateNote, "note", "", "", "更新备注")
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func removeTag(tags []string, tag string) []string {
	result := []string{}
	for _, t := range tags {
		if t != tag {
			result = append(result, t)
		}
	}
	return result
}
