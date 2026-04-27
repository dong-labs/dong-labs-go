package membercmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/member/db"
	"github.com/spf13/cobra"
)

var remindDays int

var RemindCmd = &cobra.Command{
	Use:   "remind",
	Short: "检查需要续费的会员",
	Run: func(cmd *cobra.Command, args []string) {
		days, _ := cmd.Flags().GetInt("days")
		if days == 0 {
			days = 30
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
		SELECT id, name, expire_date, member_type, project, wechat, phone, email
		FROM members
		WHERE status = 'active' AND DATE(expire_date) BETWEEN DATE(?) AND DATE(?)
		ORDER BY expire_date ASC
	`, today, targetDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name, expireDate, memberType, project, wechat, phone, email string
		rows.Scan(&id, &name, &expireDate, &memberType, &project, &wechat, &phone, &email)

		results = append(results, map[string]interface{}{
			"id":          id,
			"name":        name,
			"expire_date": expireDate,
			"member_type": memberType,
			"project":     project,
			"wechat":      wechat,
			"phone":       phone,
			"email":       email,
		})
	}

	return results, nil
}

func init() {
	RemindCmd.Flags().IntVarP(&remindDays, "days", "d", 30, "提前多少天提醒")
}
