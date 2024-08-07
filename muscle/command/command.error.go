package command

import (
	"fmt"
)

// CommandError 인터페이스 정의
type CommandError interface {
	Error() string
	Unwrap() error
}

// CommandErrorImpl 구조체 정의
type CommandErrorImpl struct {
	Message string
	Cause   error
}

// NewCommandError 함수 정의
func NewCommandError(message string, cause error) CommandError {
	return &CommandErrorImpl{
		Message: message,
		Cause:   cause,
	}
}

// Error 메서드 구현 (CommandError 인터페이스 준수)
func (c *CommandErrorImpl) Error() string {
	if c.Cause != nil {
		return fmt.Sprintf("%s: %v", c.Message, c.Cause)
	}
	return c.Message
}

// Unwrap 메서드 구현 (CommandError 인터페이스 준수)
func (c *CommandErrorImpl) Unwrap() error {
	return c.Cause
}
