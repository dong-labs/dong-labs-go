package cangcmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/cang/db"
	"github.com/spf13/cobra"
)

var searchCategory string

var SearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "搜索交易",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		category := searchCategory

		results, err := SearchTransactions(keyword, category)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(map[string]interface{}{
			"total": len(results),
			"items": results,
		})
	},
}

func SearchTransactions(keyword, category string) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	query := `SELECT id, date, amount_cents, account_id, category, note, tags, created_at, updated_at FROM transactions WHERE (note LIKE ? OR category LIKE ? OR tags LIKE ?)`
	args := []interface{}{"%" + keyword + "%", "%" + keyword + "%", "%" + keyword + "%"}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	query += " ORDER BY date DESC LIMIT 100"

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var date, accountID, cat, note, tags, createdAt, updatedAt string
		var amountCents int
		err = rows.Scan(&id, &date, &amountCents, &accountID, &cat, &note, &tags, &createdAt, &updatedAt)
		if err != nil { continue }

		results = append(results, map[string]interface{}{
			"id":         id,
			"date":       date,
			"amount":     float64(amountCents) / 100,
			"account_id":  accountID,
			"category":   cat,
			"note":       note,
			"tags":       tags,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}

	return results, nil
}

func init() {
	SearchCmd.Flags().StringVarP(&searchCategory, "category", "c", "", "按分类搜索")
}
