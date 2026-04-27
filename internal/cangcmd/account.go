package cangcmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/cang/db"
	"github.com/spf13/cobra"
)

var (
	accountName   string
	accountType   string
	accountCurrency string
)

var AccountAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "添加账户",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		actType, _ := cmd.Flags().GetString("type")
		currency, _ := cmd.Flags().GetString("currency")

		if actType == "" {
			actType = "checking"
		}
		if currency == "" {
			currency = "CNY"
		}

		database := db.GetDB()
		now := time.Now().Format("2006-01-02 15:04:05")
		result, err := database.Exec(`INSERT INTO accounts (name, type, currency, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
			name, actType, currency, now, now)
		if err != nil {
			output.PrintJSONError("INSERT_ERROR", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		output.PrintJSON(map[string]interface{}{"id": id, "name": name, "type": actType})
	},
}

func init() {
	AccountAddCmd.Flags().StringVarP(&accountType, "type", "t", "", "账户类型")
	AccountAddCmd.Flags().StringVar(&accountCurrency, "currency", "", "货币")
}
