package generate

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ShellcodeConfig struct {
	SourcePath        string            `yaml:"sourcePath"`
	IncludePath       string            `yaml:"includePath"`
	Substitutions     map[string]string `yaml:"substitutions"`
	IncludeGuard      string            `yaml:"includeGuard"`
	SourceOutputName  string            `yaml:"sourceOutputName"`
	IncludeOutputName string            `yaml:"includeOutputName"`
}

type FileConfig struct {
	SourcePath    string            `yaml:"sourcePath"`
	OutputPath    string            `yaml:"outputPath"`
	Substitutions map[string]string `yaml:"substitutions"`
}

type LoaderConfig struct {
	Token       string       `yaml:"token"`
	EncType     *string      `yaml:"enc_type"`
	KeyRequired bool         `yaml:"key_required"`
	Method      string       `yaml:"method"`
	Files       []FileConfig `yaml:"files"`
}

type Config struct {
	Shellcode ShellcodeConfig `yaml:"shellcode_template"`
	Loader    []LoaderConfig  `yaml:"loader_template"`
}

func ReadConfig(path string) (*Config, error) {
	var config Config
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
