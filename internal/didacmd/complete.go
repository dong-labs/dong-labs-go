package didacmd

import (
	"time"

	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/dida/db"
	"github.com/spf13/cobra"
)

var CompleteCmd = &cobra.Command{
	Use:   "complete <id>",
	Short: "标记待办为完成",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		database := db.GetDB()
		now := time.Now().Format("2006-01-02 15:04:05")
		_, err = database.Exec(`UPDATE todos SET status = 'done', completed_at = ?, updated_at = ? WHERE id = ?`, now, now, id)
		if err != nil {
			output.PrintJSONError("UPDATE_ERROR", err.Error())
			return
		}

		output.PrintJSON(map[string]interface{}{
			"completed": id,
			"status":    "done",
		})
	},
}
