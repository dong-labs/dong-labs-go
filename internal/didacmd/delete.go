package didacmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dida/db"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "删除待办",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		err = DeleteTodo(id)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(map[string]interface{}{
			"deleted": id,
		})
	},
}

func DeleteTodo(id int) error {
	database := db.GetDB()
	result, err := database.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return errors.NewDongError(errors.ErrInternal, "删除失败: "+err.Error(), nil)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.NewNotFoundError("Todo", id)
	}
	return nil
}
