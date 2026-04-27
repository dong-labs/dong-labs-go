package expirecmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-expire/db"
	"github.com/spf13/cobra"
)

var HistoryCmd = &cobra.Command{
	Use:   "history <id>",
	Short: "查看续费历史",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		history, err := GetRenewHistory(id)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(map[string]interface{}{
			"item_id": id,
			"total":   len(history),
			"history": history,
		})
	},
}

func GetRenewHistory(itemID int) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	rows, err := conn.Query(`
		SELECT id, old_expire_date, new_expire_date, renewed_at
		FROM renew_history
		WHERE item_id = ?
		ORDER BY renewed_at DESC
	`, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var oldDate, newDate, renewedAt string
		rows.Scan(&id, &oldDate, &newDate, &renewedAt)

		results = append(results, map[string]interface{}{
			"id":              id,
			"old_expire_date": oldDate,
			"new_expire_date": newDate,
			"renewed_at":      renewedAt,
		})
	}

	return results, nil
}
