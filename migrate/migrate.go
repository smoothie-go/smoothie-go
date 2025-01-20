package migrate

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Migrate(inputFilePath string) (string, error) {
	file, err := os.Open(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	sectionRegex := regexp.MustCompile(`^\[(.+)\]$`)
	keyValueRegex := regexp.MustCompile(`^([^:=\s][^:=]*?)\s*[:=]\s*(.*)$`) // GPT'd regex

	var output strings.Builder
	currentSection := ""
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			output.WriteString(line + "\n")
			continue
		}

		if matches := sectionRegex.FindStringSubmatch(line); matches != nil {
			currentSection = matches[1]
			output.WriteString("[" + currentSection + "]\n")
			continue
		}

		if matches := keyValueRegex.FindStringSubmatch(line); matches != nil {
			key := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])

			output.WriteString(fmt.Sprintf("%s = %s\n", key, value))
			continue
		}

		return "", fmt.Errorf("malformed line: %s", line)
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return output.String(), nil
}
