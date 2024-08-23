package generator

import (
	"fmt"
	process_error "muscle/command/error"
	generatorAnsible "muscle/generator/ansible"
	generatorGitAcitons "muscle/generator/gitactions"
	generatorTerraform "muscle/generator/terraform"
)

type Generator interface {
	CheckConfig() error
	Generate(filepath string) error
}

func NewGenerator(types string, config map[string]string) (Generator, error) {

	switch types {
	case "terraform":
		return &generatorTerraform.GeneratorTerraform{Config: config}, nil
	case "ansible":
		return &generatorAnsible.GeneratorAnsible{Config: config}, nil
	case "gitactions":
		return &generatorGitAcitons.GeneratorGitActions{Config: config}, nil
	default:
		return nil, process_error.NewError(fmt.Sprintf("invalid type '%s'", config["type"]), nil)
	}
}
