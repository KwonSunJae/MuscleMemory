package ready

type Ready interface {
	Lock() error
	LoadRepository() error
	UpdateRepository() error
	ReadyRepository() error
}

type ReadyImpl struct {
	// contains filtered or unexported fields
	Config map[string]string
}

func NewReadyProcessor(config map[string]string) (Ready, error) {
	return &ReadyImpl{Config: config}, nil
}

func (r *ReadyImpl) Lock() error {
	// Lock
	return nil
}

func (r *ReadyImpl) LoadRepository() error {
	// LoadRepository
	return nil
}

func (r *ReadyImpl) UpdateRepository() error {
	// UpdateRepository
	return nil
}

func (r *ReadyImpl) ReadyRepository() error {
	// ReadyRepository
	return nil
}
