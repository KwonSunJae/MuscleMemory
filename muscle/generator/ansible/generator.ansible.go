package generator

type GeneratorAnsible struct {
	Config map[string]string
}

func (g *GeneratorAnsible) CheckConfig() error {
	// CheckConfig
	return nil
}

func (g *GeneratorAnsible) Generate() (string, error) {
	return "", nil
}
