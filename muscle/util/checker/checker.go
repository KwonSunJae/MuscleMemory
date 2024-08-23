package checker

import (
	"fmt"
	process_error "muscle/command/error"
)

func CheckArgValidate(config map[string]string, essentialArgList []string) error {
	// Check Essential Arguments
	for _, essentialArg := range essentialArgList {
		if _, ok := config[essentialArg]; !ok {
			return process_error.NewError(fmt.Sprintf("essential argument '%s' is missing", essentialArg), nil)
		}
	}

	return nil
}
