package system

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
		return fmt.Errorf("%s %s: %v \n output message : %s ", cmdname, cmdarg, err, stdoutStderr)
	}

	c.stdout = string(stdoutStderr)

	return nil
}
