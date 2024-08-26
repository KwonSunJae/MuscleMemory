package writer

import (
	"fmt"
	"os"

	process_error "muscle/command/error"
)

func WriteConfig(filename string, configMap map[string]string) error {
	// 파일 열기
	file, err := os.Create(filename)
	if err != nil {
		return process_error.NewError("error creating file", err)
	}
	defer file.Close()

	// key=value 형식으로 파일 쓰기
	for key, value := range configMap {
		_, err := fmt.Fprintf(file, "%s=%s\n", key, value)
		if err != nil {
			return process_error.NewError("error writing file", err)
		}
	}

	return nil
}
