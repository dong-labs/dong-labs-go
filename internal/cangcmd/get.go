package cangcmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/cang/db"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "获取交易详情",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		tx, err := GetTransaction(id)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(tx)
	},
}

func GetTransaction(id int) (map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	var date, accountID, category, note, tags, createdAt, updatedAt string
	var amountCents int
	err = conn.QueryRow(`
		SELECT date, amount_cents, account_id, category, note, tags, created_at, updated_at
		FROM transactions WHERE id = ?
	`, id).Scan(&date, &amountCents, &accountID, &category, &note, &tags, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.NewNotFoundError("Transaction", id)
	}

	return map[string]interface{}{
		"id":         id,
		"date":       date,
		"amount":     float64(amountCents) / 100,
		"account_id":  accountID,
		"category":   category,
		"note":       note,
		"tags":       tags,
		"created_at": createdAt,
		"updated_at": updatedAt,
	}, nil
}
