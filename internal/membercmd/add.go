package membercmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/member/db"
	"github.com/spf13/cobra"
)

var (
	addName        string
	addWechat      string
	addPhone       string
	addEmail       string
	addAccountID   string
	addMemberType  string
	addProject     string
	addJoinDate    string
	addExpireDate  string
	addPrice       float64
	addCurrency    string
	addStatus      string
	addSource      string
	addRegion      string
	addJob         string
	addTechLevel   string
	addNotes       string
)

var AddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "添加会员",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		
		wechat, _ := cmd.Flags().GetString("wechat")
		phone, _ := cmd.Flags().GetString("phone")
		email, _ := cmd.Flags().GetString("email")
		accountID, _ := cmd.Flags().GetString("account_id")
		memberType, _ := cmd.Flags().GetString("type")
		project, _ := cmd.Flags().GetString("project")
		joinDate, _ := cmd.Flags().GetString("join_date")
		expireDate, _ := cmd.Flags().GetString("expire_date")
		price, _ := cmd.Flags().GetFloat64("price")
		currency, _ := cmd.Flags().GetString("currency")
		status, _ := cmd.Flags().GetString("status")
		source, _ := cmd.Flags().GetString("source")
		region, _ := cmd.Flags().GetString("region")
		job, _ := cmd.Flags().GetString("job")
		techLevel, _ := cmd.Flags().GetString("tech_level")
		notes, _ := cmd.Flags().GetString("notes")

		// 默认值
		if joinDate == "" {
			joinDate = time.Now().Format("2006-01-02")
		}
		if memberType == "" {
			memberType = "yearly"
		}
		if project == "" {
			project = "donglijuan"
		}
		if currency == "" {
			currency = "CNY"
		}
		if status == "" {
			status = "active"
		}

		database := db.GetDB()
		now := time.Now().Format("2006-01-02 15:04:05")
		result, err := database.Exec(`INSERT INTO members (name, wechat, phone, email, account_id, member_type, project, join_date, expire_date, price, currency, status, source, region, job, tech_level, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			name, wechat, phone, email, accountID, memberType, project, joinDate, expireDate, price, currency, status, source, region, job, techLevel, notes, now, now)
		if err != nil {
			output.PrintJSONError("INSERT_ERROR", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		output.PrintJSON(map[string]interface{}{"id": id, "name": name})
	},
}

func init() {
	AddCmd.Flags().StringVar(&addWechat, "wechat", "", "微信号")
	AddCmd.Flags().StringVar(&addPhone, "phone", "", "手机号")
	AddCmd.Flags().StringVar(&addEmail, "email", "", "邮箱")
	AddCmd.Flags().StringVar(&addAccountID, "account-id", "", "账号ID")
	AddCmd.Flags().StringVarP(&addMemberType, "type", "t", "", "会员类型")
	AddCmd.Flags().StringVarP(&addProject, "project", "p", "", "项目")
	AddCmd.Flags().StringVar(&addJoinDate, "join-date", "", "加入日期 YYYY-MM-DD")
	AddCmd.Flags().StringVar(&addExpireDate, "expire-date", "", "到期日期 YYYY-MM-DD")
	AddCmd.Flags().Float64Var(&addPrice, "price", 0, "价格")
	AddCmd.Flags().StringVar(&addCurrency, "currency", "", "货币")
	AddCmd.Flags().StringVar(&addStatus, "status", "", "状态")
	AddCmd.Flags().StringVar(&addSource, "source", "", "来源")
	AddCmd.Flags().StringVar(&addRegion, "region", "", "地区")
	AddCmd.Flags().StringVar(&addJob, "job", "", "职业")
	AddCmd.Flags().StringVar(&addTechLevel, "tech-level", "", "技术水平")
	AddCmd.Flags().StringVarP(&addNotes, "notes", "n", "", "备注")
}
