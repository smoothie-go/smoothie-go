package recipe

import (
	"fmt"
	"strings"

	"gopkg.in/ini.v1"
)

func ParseEncodingArgs(iniFilePath, inputEncArgs string) (string, error) {
	cfg, err := ini.Load(iniFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to load INI file: %w", err)
	}

	aliasToSection := make(map[string]string)
	for _, section := range cfg.Sections() {
		if section.Name() == "MACROS" || section.Name() == "" {
			continue
		}
		aliases := strings.Split(section.Name(), "/")
		for _, alias := range aliases {
			aliasToSection[alias] = section.Name()
		}
	}

	var codec string
	var result strings.Builder
	codecOptions := make(map[string]string)

	for sectionName := range aliasToSection {
		codecOptions[sectionName] = aliasToSection[sectionName]
	}

	words := strings.Fields(inputEncArgs)
	for _, word := range words {
		if result.Len() > 0 && result.String()[result.Len()-1] != ' ' {
			result.WriteString(" ")
		}

		if macroValue := cfg.Section("MACROS").Key(word).String(); macroValue != "" {
			result.WriteString(macroValue)
			continue
		}

		if codec == "" {
			if sectionName, ok := codecOptions[word]; ok {
				codec = sectionName
				continue
			}
		}

		if codec != "" {
			if presetValue := cfg.Section(codec).Key(strings.ToUpper(word)).String(); presetValue != "" {
				result.WriteString(presetValue)
				continue
			}
		}

		result.WriteString(word)
	}

	return result.String(), nil
}
