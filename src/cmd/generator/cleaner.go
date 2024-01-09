package generate

import (
	"fmt"
	"os"
	"strings"
)

func ParseArgs(arg []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, pair := range arg {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid argument: %s", pair)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		result[key] = value
	}
	return result, nil
}

func CleanShellcode(outputPath string, config *Config) {
	fmt.Println("deleting -", outputPath+"/"+config.Shellcode.SourceOutputName)
	err := os.Remove(outputPath + "/" + config.Shellcode.SourceOutputName)
	if err != nil {
		return
	}

	fmt.Println("deleting -", outputPath+"/"+config.Shellcode.IncludeOutputName)
	err = os.Remove(outputPath + "/" + config.Shellcode.IncludeOutputName)
	if err != nil {
		return
	}
}

func CleanLoader(outputPath string, config *Config) {
	for _, l := range config.Loader {
		for _, f := range l.Files {
			fmt.Println("deleting -", outputPath+"/"+f.OutputPath)
			err := os.Remove(outputPath + "/" + f.OutputPath)
			if err != nil {
				return
			}
		}
	}
}
