package init

import "fmt"

type Init interface {
	// Run the init process
	InputConfig() error
	CheckArgValidate() error
	Run() error
}

func GetInitProcessor(config map[string]string) (Init, error) {
	// set config map from cmd

	// InputConfig
	//if i.Config dont have "type" key, set default value "project"
	if _, ok := config["type"]; !ok {
		config["type"] = "default"
	}

	var tempInitProcessor Init
	switch config["type"] {
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
		return nil, fmt.Errorf("init processor error: unsupported type")
	}

	return tempInitProcessor, nil
}
