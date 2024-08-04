package generator

type GeneratorGitActions struct {
	Config map[string]string
}

func (g *GeneratorGitActions) CheckConfig() error {
	// CheckConfig
	return nil
}

func (g *GeneratorGitActions) Generate() (string, error) {
	return "", nil
}
