package Utility

import (
	"bufio"
	"os"
	"strings"
)

func GetConfig(filename string) (error, map[string]string) {
	file, err := os.Open(filename)
	if err != nil {
		return err, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	config := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[key] = value
		}
	}
	return nil, config
}