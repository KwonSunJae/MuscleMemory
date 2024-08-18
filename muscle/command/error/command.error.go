package command_error

import (
	"fmt"
	"runtime"
)

// CustomError는 사용자 정의 오류 타입입니다.
type CommandError struct {
	Message string
	File    string
	Line    int
	Cause   error
}

// Error 메소드 구현
func (e *CommandError) Error() string {
	return fmt.Sprintf("%s (file: %s, line: %d, cause: %v) ", e.Message, e.File, e.Line, e.Cause)
}

// NewCustomError는 CustomError를 생성하는 함수입니다.
func NewError(msg string, cause error) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return &CommandError{Message: msg}
	}
	return &CommandError{Message: msg, File: file, Line: line, Cause: cause}
}
