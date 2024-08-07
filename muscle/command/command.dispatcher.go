package command

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
	switch cmdType {
	case "generate-terraform":
		return c.commandHandler.GenerateTerraform(cmd)
	case "generate-ansible":
		return c.commandHandler.GenerateAnsible(cmd)
	case "init-project":
		return c.commandHandler.InitProject(cmd)
	case "ready-project":
		return c.commandHandler.ReadyProject(cmd)
	case "add-terraform-project":
		return c.commandHandler.AddTerraformProject(cmd)
	case "add-ansible-project":
		return c.commandHandler.AddAnsibleProject(cmd)
	case "enroll-git-actions":
		return c.commandHandler.EnrollGitActions(cmd)
	case "help":
		return c.commandHandler.Help(cmd)
	default:
		return "", NewCommandError("Invalid command", nil)
	}

}
