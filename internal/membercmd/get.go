package membercmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/member/db"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "获取会员详情",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		member, err := GetMember(id)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(member)
	},
}

func GetMember(id int) (map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	var name, wechat, phone, email, accountID, memberType, project, joinDate, expireDate, currency, status, source, region, job, techLevel, notes, createdAt, updatedAt string
	var price float64
	err = conn.QueryRow(`
		SELECT name, wechat, phone, email, account_id, member_type, project, join_date, expire_date, price, currency, status, source, region, job, tech_level, notes, created_at, updated_at
		FROM members WHERE id = ?
	`, id).Scan(&name, &wechat, &phone, &email, &accountID, &memberType, &project, &joinDate, &expireDate, &price, &currency, &status, &source, &region, &job, &techLevel, &notes, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.NewNotFoundError("Member", id)
	}

	return map[string]interface{}{
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
		"status":       status,
		"source":       source,
		"region":       region,
		"job":          job,
		"tech_level":   techLevel,
		"notes":        notes,
		"created_at":   createdAt,
		"updated_at":   updatedAt,
	}, nil
}
