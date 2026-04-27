package readcmd

import (
	"sort"
	"strings"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-read/db"
	"github.com/spf13/cobra"
)

type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

var TagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "列出所有标签",
	Run: func(cmd *cobra.Command, args []string) {
		tags, err := GetTags()
		if err != nil {
			printError(err)
			return
		}
		output.PrintJSON(map[string]interface{}{
			"total": len(tags),
			"tags":  tags,
		})
	},
}

func GetTags() ([]TagCount, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	rows, err := conn.Query(`SELECT tags FROM items WHERE tags != '' AND tags IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagMap := make(map[string]int)
	for rows.Next() {
		var tagsStr string
		rows.Scan(&tagsStr)
		tags := strings.Split(tagsStr, ",")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tagMap[tag]++
			}
		}
	}

	// Convert to slice and sort
	result := make([]TagCount, 0, len(tagMap))
	for tag, count := range tagMap {
		result = append(result, TagCount{Tag: tag, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result, nil
}
