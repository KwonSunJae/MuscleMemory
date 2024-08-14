package command

import (
	"fmt"
	initProcessor "muscle/command/process/init"
)

type CommandHandler interface {
	Init([]string) (string, error)
	Generate([]string) (string, error)
	Ready([]string) (string, error)
	Add([]string) (string, error)
	Enroll([]string) (string, error)
	Help([]string) (string, error)
}

type CommandHandlerImpl struct {
	// contains filtered or unexported fields
}

func NewCommandHandler() CommandHandler {
	return &CommandHandlerImpl{}
}

func (c *CommandHandlerImpl) Init(cmd []string) (string, error) {
	// Init

	// cmd to config map
	fmt.Println(cmd)
	initConfig := cmdToConfigMap(cmd)

	initProcessor, err := initProcessor.GetInitProcessor(initConfig)
	if err != nil {
		return "", fmt.Errorf("command_handler_comp_error: \n %v", err)
	}

	err = initProcessor.InputConfig()
	if err != nil {
		return "", fmt.Errorf("command_handler_input_error: \n %v", err)
	}

	err = initProcessor.CheckArgValidate()
	if err != nil {
		return "", fmt.Errorf("command_handler_validate_error: \n %v", err)
	}

	err = initProcessor.Run()
	if err != nil {
		return "", fmt.Errorf("command_handler_run_error: \n %v", err)
	}

	return "", nil
}

func (c *CommandHandlerImpl) Generate(cmd []string) (string, error) {
	// Generate
	return "", nil
}

func (c *CommandHandlerImpl) Ready(cmd []string) (string, error) {
	// Ready
	return "", nil
}

func (c *CommandHandlerImpl) Add(cmd []string) (string, error) {
	// Add
	return "", nil
}

func (c *CommandHandlerImpl) Enroll(cmd []string) (string, error) {
	// Enroll
	return "", nil
}

func (c *CommandHandlerImpl) Help(cmd []string) (string, error) {
	// Help
	return "", nil
}

// cmd to config map
func cmdToConfigMap(cmd []string) map[string]string {
	// cmd to config map
	var result = make(map[string]string)
	// if prefix is "--", set key as value and value as true
	// if prefix is "-", set key as value and value as next value
	for len(cmd) >= 0 {
		if cmd[0][:2] == "--" {
			// set key as value and value as true
			result[cmd[0][2:]] = "true"
			fmt.Println(result)
			if len(cmd) == 1 {
				break
			}
			cmd = cmd[1:]

		} else if cmd[0][:1] == "-" {
			// set key as value and value as next value
			result[cmd[0][1:]] = cmd[1]
			cmd = cmd[2:]
		}

	}

	return result
}
