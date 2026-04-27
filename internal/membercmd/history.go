package membercmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/member/db"
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
			"member_id": id,
			"total":     len(history),
			"history":   history,
		})
	},
}

func GetRenewHistory(memberID int) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	rows, err := conn.Query(`
		SELECT id, old_expire_date, new_expire_date, amount, currency, renewed_at, notes
		FROM renewals
		WHERE member_id = ?
		ORDER BY renewed_at DESC
	`, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var oldDate, newDate, currency, renewedAt, notes string
		var amount float64
		rows.Scan(&id, &oldDate, &newDate, &amount, &currency, &renewedAt, &notes)

		results = append(results, map[string]interface{}{
			"id":              id,
			"old_expire_date": oldDate,
			"new_expire_date": newDate,
			"amount":          amount,
			"currency":        currency,
			"renewed_at":       renewedAt,
			"notes":           notes,
		})
	}

	return results, nil
}
