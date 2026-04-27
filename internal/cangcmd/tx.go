package cangcmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/cang/db"
	"github.com/spf13/cobra"
)

var (
	txDate      string
	txAmount    float64
	txAccount   string
	txCategory  string
	txNote      string
	txTags      string
)

var TxCmd = &cobra.Command{
	Use:   "tx <amount>",
	Short: "记录交易",
	Run: func(cmd *cobra.Command, args []string) {
		amount, _ := cmd.Flags().GetFloat64("amount")
		date, _ := cmd.Flags().GetString("date")
		account, _ := cmd.Flags().GetString("account")
		category, _ := cmd.Flags().GetString("category")
		note, _ := cmd.Flags().GetString("note")
		tags, _ := cmd.Flags().GetString("tags")

		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		// 转换为分
		amountCents := int(amount * 100)

		database := db.GetDB()
		now := time.Now().Format("2006-01-02 15:04:05")
		result, err := database.Exec(`INSERT INTO transactions (date, amount_cents, account_id, category, note, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			date, amountCents, account, category, note, tags, now, now)
		if err != nil {
			output.PrintJSONError("INSERT_ERROR", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		output.PrintJSON(map[string]interface{}{"id": id, "amount": amount, "date": date})
	},
}

func init() {
	TxCmd.Flags().StringVarP(&txDate, "date", "d", "", "日期 YYYY-MM-DD")
	TxCmd.Flags().Float64VarP(&txAmount, "amount", "a", 0, "金额")
	TxCmd.Flags().StringVarP(&txAccount, "account", "A", "", "账户")
	TxCmd.Flags().StringVarP(&txCategory, "category", "c", "", "分类")
	TxCmd.Flags().StringVarP(&txNote, "note", "n", "", "备注")
	TxCmd.Flags().StringVarP(&txTags, "tags", "t", "", "标签")
}
