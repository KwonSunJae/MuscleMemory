package command

import (
	"fmt"
	"muscle/command/process/enroll"
	initProcessor "muscle/command/process/init"
	"muscle/command/process/ready"
	"muscle/util/crypt"
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
	Owner string
}

func NewCommandHandler() CommandHandler {
	owner, err := crypt.GetOwner()
	if err != nil {
		fmt.Println(err)
		return &CommandHandlerImpl{Owner: "NOT_FOUND"}
	}
	return &CommandHandlerImpl{Owner: owner}
}

func (c *CommandHandlerImpl) Init(cmd []string) (string, error) {
	// Init
	var types string
	// cmd to config map
	if cmd[0][0] == '-' {
		types = "default"
	} else {
		types = cmd[0]
		cmd = cmd[1:]
	}
	initConfig := cmdToConfigMap(cmd)

	initProcessor, err := initProcessor.GetInitProcessor(types, initConfig)
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

	config := make(map[string]string)
	config["project-name"] = cmd[0]
	config["work-dir"] = cmd[1]
	config["owner"] = c.Owner

	readyProcessor, err := ready.NewReadyProcessor(config)
	if err != nil {
		return "", fmt.Errorf("command_handler_comp_error: \n %v", err)
	}

	if err := readyProcessor.LoadConfig(); err != nil {
		return "", fmt.Errorf("command_handler_load_config_error: \n %v", err)
	}

	if err := readyProcessor.LoadRepository(); err != nil {
		return "", fmt.Errorf("command_handler_load_repository_error: \n %v", err)
	}

	if err := readyProcessor.Lock(); err != nil {
		return "", fmt.Errorf("command_handler_lock_error: \n %v", err)
	}

	if err := readyProcessor.ReadyRepository(); err != nil {
		return "", fmt.Errorf("command_handler_ready_repository_error: \n %v", err)
	}

	return "", nil
}

func (c *CommandHandlerImpl) Add(cmd []string) (string, error) {
	// Add
	return "", nil
}

func (c *CommandHandlerImpl) Enroll(cmd []string) (string, error) {
	// Enroll
	var conf = make(map[string]string)
	conf["project-name"] = cmd[0]
	enrollProcessor, err := enroll.NewEnrollProcessor(conf)
	if err != nil {
		return "", fmt.Errorf("command_handler_comp_error: \n %v", err)
	}

	if err := enrollProcessor.CheckLock(); err != nil {
		return "", fmt.Errorf("command_handler_check_lock_error: \n %v", err)
	}

	if err := enrollProcessor.Enroll(); err != nil {
		return "", fmt.Errorf("command_handler_enroll_error: \n %v", err)
	}

	if err := enrollProcessor.UnLock(); err != nil {
		return "", fmt.Errorf("command_handler_unlock_error: \n %v", err)
	}

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

			if len(cmd) == 1 {
				break
			}
			cmd = cmd[1:]

		} else if cmd[0][:1] == "-" {
			// set key as value and value as next value
			result[cmd[0][1:]] = cmd[1]
			if len(cmd) == 2 {
				break
			}
			cmd = cmd[2:]
		}

	}

	return result
}
