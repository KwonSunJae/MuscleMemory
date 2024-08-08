package init

type InitAnsible struct {
	Config map[string]string
}

func (i *InitAnsible) CheckConfig() error {
	// CheckConfig
	return nil
}

func (i *InitAnsible) CheckArgValidate() error {
	// CheckArgValidate
	return nil
}

func (i *InitAnsible) InputConfig() error {
	// InputConfig
	return nil
}

func (i *InitAnsible) Run() error {
	// Run
	return nil
}
