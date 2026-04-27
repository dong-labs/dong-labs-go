package membercmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/member/db"
	"github.com/dong-labs/think/internal/member/models"
	"github.com/spf13/cobra"
)

var (
	listLimit    int
	listStatus   string
	listType     string
	listProject  string
	listRegion   string
	listExpired  bool
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出会员",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		status, _ := cmd.Flags().GetString("status")
		memberType, _ := cmd.Flags().GetString("type")
		project, _ := cmd.Flags().GetString("project")
		region, _ := cmd.Flags().GetString("region")

		database := db.GetDB()
		conn, err := database.GetConnection()
		if err != nil {
			output.PrintJSONError("DB_ERROR", err.Error())
			return
		}

		query := `SELECT id, name, wechat, phone, email, account_id, member_type, project, join_date, expire_date, price, currency, status, source, region, job, tech_level, notes, created_at, updated_at FROM members WHERE 1=1`
		queryArgs := []interface{}{}

		if status != "" {
			query += " AND status = ?"
			queryArgs = append(queryArgs, status)
		}
		if memberType != "" {
			query += " AND member_type = ?"
			queryArgs = append(queryArgs, memberType)
		}
		if project != "" {
			query += " AND project = ?"
			queryArgs = append(queryArgs, project)
		}
		if region != "" {
			query += " AND region = ?"
			queryArgs = append(queryArgs, region)
		}

		query += " ORDER BY created_at DESC LIMIT ?"
		queryArgs = append(queryArgs, limit)

		rows, err := conn.Query(query, queryArgs...)
		if err != nil {
			output.PrintJSONError("QUERY_ERROR", err.Error())
			return
		}
		defer rows.Close()

		members := []models.Member{}
		for rows.Next() {
			var m models.Member
			var createdAt, updatedAt string
			var price float64
			rows.Scan(&m.ID, &m.Name, &m.Wechat, &m.Phone, &m.Email, &m.AccountID, &m.MemberType, &m.Project, &m.JoinDate, &m.ExpireDate, &price, &m.Currency, &m.Status, &m.Source, &m.Region, &m.Job, &m.TechLevel, &m.Notes, &createdAt, &updatedAt)
			m.Price = price
			m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			m.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
			members = append(members, m)
		}
		output.PrintJSON(map[string]interface{}{"total": len(members), "items": members})
	},
}

func init() {
	ListCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "返回数量")
	ListCmd.Flags().StringVarP(&listStatus, "status", "s", "", "按状态筛选")
	ListCmd.Flags().StringVarP(&listType, "type", "t", "", "按类型筛选")
	ListCmd.Flags().StringVarP(&listProject, "project", "p", "", "按项目筛选")
	ListCmd.Flags().StringVarP(&listRegion, "region", "r", "", "按地区筛选")
}
