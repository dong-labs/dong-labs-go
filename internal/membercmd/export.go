package membercmd

import (
	"encoding/json"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/member/db"
	"github.com/spf13/cobra"
)

var exportFile string

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出会员数据",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := FetchAllMembers()
		if err != nil {
			printError(err)
			return
		}

		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			output.PrintJSONError("MARSHAL_ERROR", err.Error())
			return
		}

		output.PrintJSON(map[string]interface{}{
			"exported": len(data),
			"data":     string(b),
		})
	},
}

func FetchAllMembers() ([]map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	rows, err := conn.Query(`SELECT id, name, wechat, phone, email, account_id, member_type, project, join_date, expire_date, price, currency, status, source, region, job, tech_level, notes, created_at, updated_at FROM members ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name, wechat, phone, email, accountID, memberType, project, joinDate, expireDate, currency, status, source, region, job, techLevel, notes, createdAt, updatedAt string
		var price float64
		err = rows.Scan(&id, &name, &wechat, &phone, &email, &accountID, &memberType, &project, &joinDate, &expireDate, &price, &currency, &status, &source, &region, &job, &techLevel, &notes, &createdAt, &updatedAt)
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
			"status":       status,
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
	ExportCmd.Flags().StringVarP(&exportFile, "output", "o", "members.json", "输出文件")
}
