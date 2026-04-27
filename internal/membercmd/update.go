package membercmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/member/db"
	"github.com/spf13/cobra"
)

var (
	updateName        string
	updateWechat      string
	updatePhone       string
	updateEmail       string
	updateAccountID   string
	updateMemberType  string
	updateProject     string
	updateJoinDate    string
	updateExpireDate  string
	updatePrice       float64
	updateCurrency    string
	updateStatus      string
	updateSource      string
	updateRegion      string
	updateJob         string
	updateTechLevel   string
	updateNotes       string
)

var UpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "更新会员",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		updates := make(map[string]interface{})
		if cmd.Flags().Changed("name") { updates["name"] = updateName }
		if cmd.Flags().Changed("wechat") { updates["wechat"] = updateWechat }
		if cmd.Flags().Changed("phone") { updates["phone"] = updatePhone }
		if cmd.Flags().Changed("email") { updates["email"] = updateEmail }
		if cmd.Flags().Changed("account-id") { updates["account_id"] = updateAccountID }
		if cmd.Flags().Changed("type") { updates["member_type"] = updateMemberType }
		if cmd.Flags().Changed("project") { updates["project"] = updateProject }
		if cmd.Flags().Changed("join-date") { updates["join_date"] = updateJoinDate }
		if cmd.Flags().Changed("expire-date") { updates["expire_date"] = updateExpireDate }
		if cmd.Flags().Changed("price") { updates["price"] = updatePrice }
		if cmd.Flags().Changed("currency") { updates["currency"] = updateCurrency }
		if cmd.Flags().Changed("status") { updates["status"] = updateStatus }
		if cmd.Flags().Changed("source") { updates["source"] = updateSource }
		if cmd.Flags().Changed("region") { updates["region"] = updateRegion }
		if cmd.Flags().Changed("job") { updates["job"] = updateJob }
		if cmd.Flags().Changed("tech-level") { updates["tech_level"] = updateTechLevel }
		if cmd.Flags().Changed("notes") { updates["notes"] = updateNotes }
		updates["updated_at"] = time.Now().Format("2006-01-02 15:04:05")

		if len(updates) == 1 {
			output.PrintJSONError("VALIDATION_ERROR", "没有指定要更新的字段")
			return
		}

		result, err := UpdateMember(id, updates)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(result)
	},
}

func UpdateMember(id int, updates map[string]interface{}) (map[string]interface{}, error) {
	database := db.GetDB()
	setClause := ""
	args := make([]interface{}, 0)
	for key, val := range updates {
		if setClause != "" { setClause += ", " }
		setClause += key + " = ?"
		args = append(args, val)
	}
	args = append(args, id)
	_, err := database.Exec("UPDATE members SET "+setClause+" WHERE id = ?", args...)
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "更新失败: "+err.Error(), nil)
	}
	return GetMember(id)
}

func init() {
	UpdateCmd.Flags().StringVar(&updateName, "name", "", "姓名")
	UpdateCmd.Flags().StringVar(&updateWechat, "wechat", "", "微信号")
	UpdateCmd.Flags().StringVar(&updatePhone, "phone", "", "手机号")
	UpdateCmd.Flags().StringVar(&updateEmail, "email", "", "邮箱")
	UpdateCmd.Flags().StringVar(&updateAccountID, "account-id", "", "账号ID")
	UpdateCmd.Flags().StringVarP(&updateMemberType, "type", "t", "", "会员类型")
	UpdateCmd.Flags().StringVarP(&updateProject, "project", "p", "", "项目")
	UpdateCmd.Flags().StringVar(&updateJoinDate, "join-date", "", "加入日期")
	UpdateCmd.Flags().StringVar(&updateExpireDate, "expire-date", "", "到期日期")
	UpdateCmd.Flags().Float64Var(&updatePrice, "price", 0, "价格")
	UpdateCmd.Flags().StringVar(&updateCurrency, "currency", "", "货币")
	UpdateCmd.Flags().StringVar(&updateStatus, "status", "", "状态")
	UpdateCmd.Flags().StringVar(&updateSource, "source", "", "来源")
	UpdateCmd.Flags().StringVar(&updateRegion, "region", "", "地区")
	UpdateCmd.Flags().StringVar(&updateJob, "job", "", "职业")
	UpdateCmd.Flags().StringVar(&updateTechLevel, "tech-level", "", "技术水平")
	UpdateCmd.Flags().StringVarP(&updateNotes, "notes", "n", "", "备注")
}
