package membercmd

import (
	"strconv"

	"github.com/dong-labs/think/internal/core/errors"
	"github.com/dong-labs/think/internal/core/output"
)

func parseID(s string) (int, error) {
	return strconv.Atoi(s)
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
