package cangcmd

import (
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/cang/db"
	"github.com/spf13/cobra"
)

var listLimit int

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出交易记录",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")

		database := db.GetDB()
		conn, err := database.GetConnection()
		if err != nil {
			output.PrintJSONError("DB_ERROR", err.Error())
			return
		}

		rows, err := conn.Query(`SELECT id, date, amount_cents, account_id, category, note, tags, created_at, updated_at FROM transactions ORDER BY date DESC LIMIT ?`, limit)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		results := []map[string]interface{}{}
		for rows.Next() {
			var id int
			var date, accountID, category, note, tags, createdAt, updatedAt string
			var amountCents int
			rows.Scan(&id, &date, &amountCents, &accountID, &category, &note, &tags, &createdAt, &updatedAt)
			
			results = append(results, map[string]interface{}{
				"id":         id,
				"date":       date,
				"amount":     float64(amountCents) / 100,
				"account_id":  accountID,
				"category":   category,
				"note":       note,
				"tags":       tags,
				"created_at": createdAt,
				"updated_at": updatedAt,
			})
		}
		output.PrintJSON(map[string]interface{}{"total": len(results), "items": results})
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "返回数量")
}
