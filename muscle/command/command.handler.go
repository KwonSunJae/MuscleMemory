package command

type CommandHandler interface {
	Init([]string) (string, error)
	GenerateTerraform([]string) (string, error)
	GenerateAnsible([]string) (string, error)
	InitProject([]string) (string, error)
	ReadyProject([]string) (string, error)
	AddTerraformProject([]string) (string, error)
	AddAnsibleProject([]string) (string, error)
	EnrollGitActions([]string) (string, error)
	Help([]string) (string, error)
}

type CommandHandlerImpl struct {
	// contains filtered or unexported fields
}

func NewCommandHandler() CommandHandler {
	return &CommandHandlerImpl{}
}

func (c *CommandHandlerImpl) GenerateTerraform(cmd []string) (string, error) {
	// GenerateTerraform

	return "", nil
}

func (c *CommandHandlerImpl) GenerateAnsible(cmd []string) (string, error) {
	// GenerateAnsible

	return "", nil
}

func (c *CommandHandlerImpl) InitProject(cmd []string) (string, error) {
	// InitProject

	// 0 : project type (terraform, ansible),  1: project name
	var projectType = cmd[0]
	var projectName = cmd[1]
	// 0. Check Project Name is valid and not duplicated

	// 1. Create a new directory

	// 2. Create a new git repository

	// 3. Create a new terraform project

	return "", nil
}

func (c *CommandHandlerImpl) ReadyProject(cmd []string) (string, error) {
	// ReadyProject

	return "", nil
}

func (c *CommandHandlerImpl) AddTerraformProject(cmd []string) (string, error) {
	// AddTerraformProject

	return "", nil
}

func (c *CommandHandlerImpl) AddAnsibleProject(cmd []string) (string, error) {
	// AddAnsibleProject

	return "", nil
}

func (c *CommandHandlerImpl) EnrollGitActions(cmd []string) (string, error) {
	// EnrollGitActions

	return "", nil
}

func (c *CommandHandlerImpl) Help(cmd []string) (string, error) {
	// Help

	return "", nil
}

func (c *CommandHandlerImpl) Init(cmd []string) (string, error) {
	// Init

	return "", nil
}
