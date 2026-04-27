package expirecmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-expire/db"
	"github.com/spf13/cobra"
)

var renewNewDate string

var RenewCmd = &cobra.Command{
	Use:   "renew <id>",
	Short: "续费订阅",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		newDate, _ := cmd.Flags().GetString("new_date")

		result, err := RenewItem(id, newDate)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(result)
	},
}

func RenewItem(id int, newDate string) (map[string]interface{}, error) {
	database := db.GetDB()

	// Get current item
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	var oldDate string
	err = conn.QueryRow("SELECT expire_date FROM items WHERE id = ?", id).Scan(&oldDate)
	if err != nil {
		return nil, errors.NewNotFoundError("Item", id)
	}

	// If no new date specified, try to calculate from old date
	if newDate == "" {
		// Default: add 1 year
		parsedDate, _ := time.Parse("2006-01-02", oldDate)
		newDate = parsedDate.AddDate(1, 0, 0).Format("2006-01-02")
	}

	// Update item
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err = database.Exec("UPDATE items SET expire_date = ?, updated_at = ? WHERE id = ?", newDate, now, id)
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "更新失败: "+err.Error(), nil)
	}

	// Record history
	_, err = database.Exec(`INSERT INTO renew_history (item_id, old_expire_date, new_expire_date, renewed_at) VALUES (?, ?, ?, ?)`,
		id, oldDate, newDate, now)

	// Return updated item
	return GetItem(id)
}

func init() {
	RenewCmd.Flags().StringVarP(&renewNewDate, "new_date", "n", "", "新到期日期 YYYY-MM-DD")
}
