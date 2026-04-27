package didacmd

import (
	"strconv"
	"strings"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
)

func parseID(s string) (int, error) {
	return strconv.Atoi(s)
}

func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	tags := strings.Split(tagsStr, ",")
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func printError(err error) {
	if e, ok := err.(*errors.DongError); ok {
		output.PrintJSONError(string(e.Code), e.Message)
	} else if e, ok := err.(*errors.NotFoundError); ok {
		output.PrintJSONError(string(errors.ErrNotFound), e.Error())
	} else if e, ok := err.(*errors.ValidationError); ok {
		output.PrintJSONError(string(e.Code), e.Error())
	} else {
		output.PrintJSONError("ERROR", err.Error())
	}
}
