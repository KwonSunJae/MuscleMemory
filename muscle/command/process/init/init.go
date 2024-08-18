package init

import (
	"bufio"
	"fmt"
	process_error "muscle/command/error"
	"os"
	"strings"
)

type Init interface {
	// Run the init process
	InputConfig() error
	CheckArgValidate() error
	Run() error
}

func GetInitProcessor(config map[string]string) (Init, error) {
	// set config map from cmd

	// InputConfig
	//if i.Config dont have "type" key, set default value "project"

	var tempInitProcessor Init
	switch config["type"] {
	case "project":
		tempInitProcessor = &InitProject{Config: config}
	case "terraform":
		tempInitProcessor = &InitTerraform{Config: config}
	case "ansible":
		tempInitProcessor = &InitAnsible{Config: config}
	case "gitactions":
		tempInitProcessor = &InitGitActions{Config: config}
	case "default":
		tempInitProcessor = &InitDefault{Config: config}
	default:
		return nil, process_error.NewError(fmt.Sprintf("invalid type '%s'", config["type"]), nil)
	}

	return tempInitProcessor, nil
}

// LoadConfig 함수는 주어진 파일에서 key=value 쌍을 읽어 map에 저장합니다.
func LoadConfig(filename string) (map[string]string, error) {
	configMap := make(map[string]string)

	// 파일 열기
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// 파일을 한 줄씩 읽기
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 주석 처리된 줄 생략 (첫 글자가 #인 경우)
		if strings.HasPrefix(line, "#") {
			continue
		}

		// key=value 형식으로 분리
		parts := strings.SplitN(line, "=", 2) // 최대 2개의 부분으로 분리
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])   // key의 앞뒤 공백 제거
			value := strings.TrimSpace(parts[1]) // value의 앞뒤 공백 제거
			configMap[key] = value
		}
	}

	// 에러 체크
	if err := scanner.Err(); err != nil {
		return nil, process_error.NewError("error reading file", err)
	}

	return configMap, nil
}

func WriteConfig(filename string, configMap map[string]string) error {
	// 파일 열기
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// key=value 형식으로 파일 쓰기
	for key, value := range configMap {
		_, err := fmt.Fprintf(file, "%s=%s\n", key, value)
		if err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}
	}

	return nil
}

func CheckArgValidate(config map[string]string, essentialArgList []string) error {
	// Check Essential Arguments
	for _, essentialArg := range essentialArgList {
		if _, ok := config[essentialArg]; !ok {
			return process_error.NewError(fmt.Sprintf("essential argument '%s' is missing", essentialArg), nil)
		}
	}

	return nil
}
