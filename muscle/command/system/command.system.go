package command

import (
	"fmt"
	"os/exec"
)

type CommandSystem interface {
	Execute(cmdname string, cmdarg ...string) error
}

type CommandSystemExecutor struct {
	// contains filtered or unexported fields
	stdout string
	stderr string
}

func NewCommandSystemExecutor() *CommandSystemExecutor {
	return &CommandSystemExecutor{}
}

func (c *CommandSystemExecutor) Execute(cmdname string, cmdarg ...string) error {
	// Execute

	cmd := exec.Command(cmdname, cmdarg...)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		c.stderr = string(stdoutStderr)
		return fmt.Errorf("command_system_executor_error:%v", err)
	}

	c.stdout = string(stdoutStderr)

	return nil
}
