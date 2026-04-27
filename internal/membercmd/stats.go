package membercmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/member/db"
	"github.com/spf13/cobra"
)

var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "统计信息",
	Run: func(cmd *cobra.Command, args []string) {
		stats, err := GetStats()
		if err != nil {
			printError(err)
			return
		}
		output.PrintJSON(stats)
	},
}

type StatsResponse struct {
	Total     int              `json:"total"`
	Active    int              `json:"active"`
	Expired   int              `json:"expired"`
	ThisMonth int              `json:"this_month"`
	ThisYear  int              `json:"this_year"`
	ByType    []TypeStat       `json:"by_type,omitempty"`
	ByProject []ProjectStat    `json:"by_project,omitempty"`
}

type TypeStat struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type ProjectStat struct {
	Project string `json:"project"`
	Count   int    `json:"count"`
}

func GetStats() (*StatsResponse, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	stats := &StatsResponse{}
	today := time.Now().Format("2006-01-02")

	// Total
	conn.QueryRow("SELECT COUNT(*) FROM members").Scan(&stats.Total)

	// Active
	conn.QueryRow("SELECT COUNT(*) FROM members WHERE status = 'active'").Scan(&stats.Active)

	// Expired
	conn.QueryRow("SELECT COUNT(*) FROM members WHERE expire_date < ?", today).Scan(&stats.Expired)

	// This month
	monthStart, _ := time.Parse("2006-01-02", today[:8]+"01")
	conn.QueryRow("SELECT COUNT(*) FROM members WHERE join_date >= ?", monthStart.Format("2006-01-02")).Scan(&stats.ThisMonth)

	// This year
	yearStart := today[:5] + "01-01"
	conn.QueryRow("SELECT COUNT(*) FROM members WHERE join_date >= ?", yearStart).Scan(&stats.ThisYear)

	// By type
	rows, err := conn.Query(`SELECT member_type, COUNT(*) as count FROM members GROUP BY member_type ORDER BY count DESC`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t string
			var c int
			rows.Scan(&t, &c)
			stats.ByType = append(stats.ByType, TypeStat{Type: t, Count: c})
		}
	}

	// By project
	rows2, _ := conn.Query(`SELECT project, COUNT(*) as count FROM members GROUP BY project ORDER BY count DESC`)
	if rows2 != nil {
		defer rows2.Close()
		for rows2.Next() {
			var p string
			var c int
			rows2.Scan(&p, &c)
			stats.ByProject = append(stats.ByProject, ProjectStat{Project: p, Count: c})
		}
	}

	return stats, nil
}
