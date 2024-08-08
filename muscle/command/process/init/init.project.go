package init

type InitProject struct {
	Config map[string]string
}

func (i *InitProject) CheckConfig() error {
	// CheckConfig
	return nil
}

func (i *InitProject) CheckArgValidate() error {
	// CheckArgValidate
	return nil
}

func (i *InitProject) InputConfig() error {
	// InputConfig
	return nil
}

func (i *InitProject) Run() error {
	// Run
	return nil
}
