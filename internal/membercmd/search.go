package membercmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/member/db"
	"github.com/spf13/cobra"
)

var searchStatus, searchType string

var SearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "搜索会员",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		status := searchStatus
		memberType := searchType

		results, err := SearchMembers(keyword, status, memberType)
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

func SearchMembers(keyword, status, memberType string) ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	query := `SELECT id, name, wechat, phone, email, account_id, member_type, project, join_date, expire_date, price, currency, status, source, region, job, tech_level, notes, created_at, updated_at FROM members WHERE (name LIKE ? OR wechat LIKE ? OR phone LIKE ? OR email LIKE ? OR notes LIKE ?)`
	args := []interface{}{"%" + keyword + "%", "%" + keyword + "%", "%" + keyword + "%", "%" + keyword + "%", "%" + keyword + "%"}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if memberType != "" {
		query += " AND member_type = ?"
		args = append(args, memberType)
	}

	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name, wechat, phone, email, accountID, memberType, project, joinDate, expireDate, currency, stat, source, region, job, techLevel, notes, createdAt, updatedAt string
		var price float64
		err = rows.Scan(&id, &name, &wechat, &phone, &email, &accountID, &memberType, &project, &joinDate, &expireDate, &price, &currency, &stat, &source, &region, &job, &techLevel, &notes, &createdAt, &updatedAt)
		if err != nil { continue }

		results = append(results, map[string]interface{}{
			"id":           id,
			"name":         name,
			"wechat":       wechat,
			"phone":        phone,
			"email":        email,
			"account_id":   accountID,
			"member_type":  memberType,
			"project":      project,
			"join_date":    joinDate,
			"expire_date":  expireDate,
			"price":        price,
			"currency":     currency,
			"status":       stat,
			"source":       source,
			"region":       region,
			"job":          job,
			"tech_level":   techLevel,
			"notes":        notes,
			"created_at":   createdAt,
			"updated_at":   updatedAt,
		})
	}

	return results, nil
}

func init() {
	SearchCmd.Flags().StringVarP(&searchStatus, "status", "s", "", "按状态搜索")
	SearchCmd.Flags().StringVarP(&searchType, "type", "t", "", "按类型搜索")
}
