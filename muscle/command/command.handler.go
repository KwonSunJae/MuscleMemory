package command

import (
	"fmt"
	initProcessor "muscle/command/process/init"
	"muscle/command/process/ready"
	"muscle/util/crypt"
	"net"
	"os"
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
	owner, err := GetOwner()
	if err != nil {
		return &CommandHandlerImpl{Owner: "NOT_FOUND"}
	}
	return &CommandHandlerImpl{Owner: owner}
}

func GetOwner() (string, error) {
	userName, err := getUserName()
	if err != nil {
		return "", err
	}

	macAddr, err := getMACAddress()
	if err != nil {
		return "", err
	}

	return crypt.Encrypt(userName, macAddr)
}

func getMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		if iface.HardwareAddr != nil {
			return iface.HardwareAddr.String(), nil
		}
	}
	return "", fmt.Errorf("MAC 주소를 찾을 수 없습니다.")
}

func getUserName() (string, error) {
	userName := os.Getenv("USER") // Unix 계열 (Linux, macOS)
	if userName == "" {
		userName = os.Getenv("USERNAME") // Windows
	}

	if userName == "" {
		return "", fmt.Errorf("사용자 이름을 찾을 수 없습니다.")
	}

	return userName, nil
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
