package passcmd

import (
	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
	"github.com/dong-labs/think/internal/pass/db"
	"github.com/spf13/cobra"
)

var showPassword bool

var GetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "获取密码详情",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			output.PrintJSONError("PARSE_ERROR", err.Error())
			return
		}

		showPass, _ := cmd.Flags().GetBool("password")
		
		password, err := GetPassword(id, showPass)
		if err != nil {
			printError(err)
			return
		}

		output.PrintJSON(password)
	},
}

func GetPassword(id int, showPassword bool) (map[string]interface{}, error) {
	database := db.GetDB()
	conn, err := database.GetConnection()
	if err != nil {
		return nil, errors.NewDongError(errors.ErrInternal, "获取数据库连接失败: "+err.Error(), nil)
	}

	var title, username, password, url, category, tags, notes, createdAt, updatedAt string
	err = conn.QueryRow(`
		SELECT title, username, password, url, category, tags, notes, created_at, updated_at
		FROM passwords WHERE id = ?
	`, id).Scan(&title, &username, &password, &url, &category, &tags, &notes, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.NewNotFoundError("Password", id)
	}

	result := map[string]interface{}{
		"id":         id,
		"title":      title,
		"username":   username,
		"url":        url,
		"category":   category,
		"tags":       tags,
		"notes":      notes,
		"created_at": createdAt,
		"updated_at": updatedAt,
	}
	
	if showPassword {
		result["password"] = password
	}
	
	return result, nil
}

func init() {
	GetCmd.Flags().BoolVarP(&showPassword, "password", "p", false, "显示密码")
}
