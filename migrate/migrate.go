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
	keyValueRegex := regexp.MustCompile(`^([^:=\s][^:=]*?)\s*[:=]\s*(.*)$`) // Matches key-value pairs

	yes := map[string]bool{
		"on": true, "true": true, "yes": true, "y": true, "1": true,
		"yeah": true, "yea": true, "yep": true, "sure": true, "positive": true,
	}
	no := map[string]bool{
		"off": true, "false": true, "no": true, "n": true, "nah": true,
		"nope": true, "negative": true, "negatory": true, "0": true, "0.0": true,
		"null": true, "": true, " ": true, "\t": true, "none": true,
	}

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

			lowerValue := strings.ToLower(value)
			if yes[lowerValue] {
				value = "true"
			} else if no[lowerValue] {
				value = "false"
			}

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
