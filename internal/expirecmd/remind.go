package expirecmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dong-expire/db"
	"github.com/spf13/cobra"
)

var remindDays int

var RemindCmd = &cobra.Command{
	Use:   "remind",
	Short: "检查需要提醒的订阅",
	Run: func(cmd *cobra.Command, args []string) {
		days, _ := cmd.Flags().GetInt("days")
		if days == 0 {
			days = 7
		}

		items, err := GetReminders(days)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(map[string]interface{}{
			"days":  days,
			"total": len(items),
			"items": items,
		})
	},
}

func GetReminders(days int) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	today := time.Now().Format("2006-01-02")
	targetDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")

	rows, err := conn.Query(`
		SELECT id, name, category, expire_date, reminder_days
		FROM items
		WHERE DATE(expire_date) BETWEEN DATE(?) AND DATE(?)
		ORDER BY expire_date ASC
	`, today, targetDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name, category, expireDate, reminderDays string
		rows.Scan(&id, &name, &category, &expireDate, &reminderDays)

		// Calculate days until expiration
		expiredTime, _ := time.Parse("2006-01-02", expireDate)
		daysUntil := int(time.Until(expiredTime).Hours() / 24)

		results = append(results, map[string]interface{}{
			"id":            id,
			"name":          name,
			"category":      category,
			"expire_date":   expireDate,
			"reminder_days": reminderDays,
			"days_until":    daysUntil,
		})
	}

	return results, nil
}

func init() {
	RemindCmd.Flags().IntVarP(&remindDays, "days", "d", 7, "提前多少天提醒")
}
