package init

import (
	"fmt"
	process_error "muscle/command/error"
)

type Init interface {
	// Run the init process
	InputConfig() error
	CheckArgValidate() error
	Run() error
}

func GetInitProcessor(types string, config map[string]string) (Init, error) {
	// set config map from cmd

	// InputConfig
	//if i.Config dont have "type" key, set default value "project"

	var tempInitProcessor Init
	switch types {
	case "project":
		tempInitProcessor = &InitProject{Config: config}
	case "terraform":
		tempInitProcessor = &InitTerraform{Config: config}
	case "ansible":
		tempInitProcessor = &InitAnsible{Config: config}
	case "gitactions":
		tempInitProcessor = &InitGitActions{Config: config}
	case "default":
		tempInitProcessor = &InitDefault{Config: config}
	default:
		return nil, process_error.NewError(fmt.Sprintf("invalid type '%s'", config["type"]), nil)
	}

	return tempInitProcessor, nil
}
