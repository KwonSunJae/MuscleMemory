package loader

import (
	"bufio"
	"fmt"
	process_error "muscle/command/error"
	"os"
	"strings"
)

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
