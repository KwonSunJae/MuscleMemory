package generator

type GeneratorAnsible struct {
	Config map[string]string
}

func (g *GeneratorAnsible) CheckConfig() error {
	// CheckConfig
	return nil
}

func (g *GeneratorAnsible) Generate(filepath string) error {
	return nil
}
