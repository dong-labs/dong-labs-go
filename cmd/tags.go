// Package cmd 提供 tags 命令
package cmd

import (
	"strings"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/think/db"
	"github.com/dong-labs/think/internal/think/models"
	"github.com/spf13/cobra"
)

// tagsCmd tags 命令
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "列出所有标签",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.GetDB()

		rows, err := database.Query(`
			SELECT tags, COUNT(*) as count
			FROM thoughts
			WHERE tags IS NOT NULL AND tags != ''
			GROUP BY tags
			ORDER BY count DESC
		`)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		tagMap := make(map[string]int)
		var totalTags int

		for rows.Next() {
			var tagsStr string
			var count int
			rows.Scan(&tagsStr, &count)

			// 分割标签（用逗号分隔）
			tagList := strings.Split(tagsStr, ",")
			for _, tag := range tagList {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					tagMap[tag] += count
					totalTags++
				}
			}
		}

		// 转换为列表
		tags := make([]models.TagStat, 0, len(tagMap))
		for tag, count := range tagMap {
			tags = append(tags, models.TagStat{Tag: tag, Count: count})
		}

		// 按数量排序
		for i := 0; i < len(tags); i++ {
			for j := i + 1; j < len(tags); j++ {
				if tags[j].Count > tags[i].Count {
					tags[i], tags[j] = tags[j], tags[i]
				}
			}
		}

		output.PrintJSON(models.TagsResponse{
			Total: len(tags),
			Tags:  tags,
		})
	},
}
