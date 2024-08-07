package generator

import (
	"errors"
	"fmt"
	"io"
	generatorAnsible "muscle/generator/ansible"
	generatorGitAcitons "muscle/generator/gitactions"
	generatorTerraform "muscle/generator/terraform"
	"os"
)

type Generator interface {
	CheckConfig() error
	Generate() (string, error)
}

type GeneratorBase struct {
	Config map[string]string
}

func (g *GeneratorBase) NewGeneratorWithFile(filepath string, types string) (Generator, error) {
	//file read
	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.New("file path is invalid : " + filepath)
	}

	//readline file parse to map
	g.Config = make(map[string]string)
	//parse file to map
	for {
		var key, value string
		_, err := fmt.Fscanf(file, "%s=%s\n", &key, &value)
		//if key containes #, skip
		if key[0] == '#' {
			continue
		}
		//if EOF, break
		if err == io.EOF {
			break
		}
		g.Config[key] = value
	}

	switch types {
	case "terraform":
		return &generatorTerraform.GeneratorTerraform{Config: g.Config}, nil
	case "ansible":
		return &generatorAnsible.GeneratorAnsible{Config: g.Config}, nil
	case "gitactions":
		return &generatorGitAcitons.GeneratorGitActions{Config: g.Config}, nil
	default:
		return nil, errors.New("unsupported type")
	}
}
