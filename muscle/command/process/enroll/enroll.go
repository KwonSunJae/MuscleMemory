package enroll

type Enroll interface {
	// Run the enroll process
	CheckLock() error
	Enroll() error
	UnLock() error
}

type EnrollImpl struct {
	// contains filtered or unexported fields
	Config map[string]string
}

func NewEnrollProcessor(config map[string]string) (Enroll, error) {
	return &EnrollImpl{Config: config}, nil
}

func (e *EnrollImpl) CheckLock() error {
	// CheckLock
	return nil
}

func (e *EnrollImpl) Enroll() error {
	// Enroll
	return nil
}

func (e *EnrollImpl) UnLock() error {
	// UnLock
	return nil
}
