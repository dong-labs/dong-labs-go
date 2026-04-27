package timelinecmd

import (
	"strings"
	"time"

	"github.com/dong-labs/think/internal/core/dates"
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/timeline/db"
	"github.com/spf13/cobra"
)

var (
	updateTitle       string
	updateDescription string
	updateCategory    string
	updateTags        []string
	updateDate        string
)

var UpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "更新事件",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		data := make(map[string]interface{})
		if cmd.Flags().Changed("title") {
			data["title"] = updateTitle
		}
		if cmd.Flags().Changed("description") {
			data["description"] = updateDescription
		}
		if cmd.Flags().Changed("category") {
			data["category"] = updateCategory
		}
		if cmd.Flags().Changed("tags") {
			data["tags"] = strings.Join(updateTags, ",")
		}
		if cmd.Flags().Changed("date") {
			date, err := dates.Parse(updateDate)
			if err != nil {
				output.PrintJSONError("PARSE_ERROR", err.Error())
				return
			}
			data["date"] = date
		}
		data["updated_at"] = time.Now().Format("2006-01-02 15:04:05")

		if len(data) == 1 { // Only updated_at
			output.PrintJSONError("VALIDATION_ERROR", "没有指定要更新的字段")
			return
		}

		result, err := UpdateEvent(id, data)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(result)
	},
}

func UpdateEvent(id int, data map[string]interface{}) (map[string]interface{}, error) {
	database := db.GetDB()

	// Build UPDATE query
	setClause := ""
	args := make([]interface{}, 0)
	for key, val := range data {
		if setClause != "" {
			setClause += ", "
		}
		setClause += key + " = ?"
		args = append(args, val)
	}
	args = append(args, id)

	_, err := database.Exec("UPDATE events SET "+setClause+" WHERE id = ?", args...)
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "更新失败: "+err.Error(), nil)
	}

	return GetEvent(id)
}

func init() {
	UpdateCmd.Flags().StringVarP(&updateTitle, "title", "t", "", "事件标题")
	UpdateCmd.Flags().StringVarP(&updateDescription, "description", "d", "", "事件描述")
	UpdateCmd.Flags().StringVarP(&updateCategory, "category", "c", "", "事件分类")
	UpdateCmd.Flags().StringSliceVar(&updateTags, "tags", []string{}, "标签")
	UpdateCmd.Flags().StringVar(&updateDate, "date", "", "事件日期")
}
