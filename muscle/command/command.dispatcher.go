package command

import "fmt"

type CommandDispatcher interface {
	CommandDispatch(cmd []string) (string, error)
}

type CommandDispatcherImpl struct {
	// contains filtered or unexported fields
	commandHandler CommandHandler
}

func NewCommandDispatcher() CommandDispatcher {
	handler := NewCommandHandler()

	return &CommandDispatcherImpl{commandHandler: handler}
}

func (c *CommandDispatcherImpl) CommandDispatch(cmd []string) (string, error) {
	// CommandDispatch

	cmdType := cmd[0]
	cmd = cmd[1:]
	var out string
	var err error
	switch cmdType {
	case "init":
		out, err = c.commandHandler.Init(cmd)
	case "generate":
		out, err = c.commandHandler.Generate(cmd)
	case "ready":
		out, err = c.commandHandler.Ready(cmd)
	case "add":
		out, err = c.commandHandler.Add(cmd)
	case "enroll":
		out, err = c.commandHandler.Enroll(cmd)
	case "help":
		out, err = c.commandHandler.Help(cmd)
	default:
		out, err = "", fmt.Errorf("command_dispatcher_error: unsupported command")
	}

	if err != nil {
		return out, NewCommandError("command_dispatcher_error", err)
	}

	return out, nil

}
