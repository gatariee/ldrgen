package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"ldrgen/cmd/utils"
)

// Config and other custom types would be defined elsewhere in the package.

func ProcessShellcodeTemplate(binPath, enc string, args map[string]string, outputPath string, config *Config, templatePath string) error {
	if enc != "" && args["key"] == "" {
		return fmt.Errorf("key not provided")
	}

	if enc == "xor" {
		if err := encryptShellcodeXOR(binPath, args["key"]); err != nil {
			return err
		}
		binPath += ".enc"
	}

	shellcodeArray, err := ToCArray(binPath)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(binPath)
	if err != nil {
		return err
	}

	shellcodeTemplate, err := prepareShellcodeTemplate(templatePath, config, shellcodeArray, fileInfo.Size())
	if err != nil {
		return err
	}

	outputFilePath, err := writeShellcodeTemplate(outputPath, config.Shellcode.IncludeOutputName, shellcodeTemplate)
	if err != nil {
		return err
	}

	printSummary(fileInfo.Size(), shellcodeArray, outputFilePath, config.Shellcode.IncludeOutputName)

	return nil
}

func encryptShellcodeXOR(binPath, key string) error {
	shellcode, err := os.ReadFile(binPath)
	if err != nil {
		return err
	}

	printEncryptionStart("XOR", key, shellcode)

	for i := range shellcode {
		shellcode[i] ^= key[i%len(key)]
	}

	printEncryptionEnd(shellcode)

	return os.WriteFile(binPath+".enc", shellcode, 0o644)
}

func prepareShellcodeTemplate(templatePath string, config *Config, shellcodeArray string, fileSize int64) (string, error) {
	shellcodeTemplateContent, err := ioutil.ReadFile(filepath.Join(templatePath, config.Shellcode.IncludePath))
	if err != nil {
		return "", err
	}

	shellcodeTemplate := strings.ReplaceAll(string(shellcodeTemplateContent), config.Shellcode.Substitutions["shellcode"], shellcodeArray)
	shellcodeTemplate = strings.ReplaceAll(shellcodeTemplate, config.Shellcode.Substitutions["shellcode_size"], fmt.Sprintf("%d", fileSize))

	return shellcodeTemplate, nil
}

func writeShellcodeTemplate(outputPath, outputName, content string) (string, error) {
	fullOutputPath, err := filepath.Abs(filepath.Join(strings.TrimSpace(outputPath), outputName))
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(fullOutputPath, []byte(content), 0o644); err != nil {
		return "", err
	}

	return fullOutputPath, nil
}

func printSummary(fileSize int64, shellcodeArray, outputFilePath, includeOutputName string) {
	utils.Print(fmt.Sprintf("Size: %s bytes", color.New(color.Bold).Sprintf("%d", fileSize)), true)
	utils.Print(fmt.Sprintf("Starting bytes: %s ... }", color.New(color.Bold).Sprintf("%s", shellcodeArray[:24])), true)
	utils.PrintWhite(fmt.Sprintf("[ %s ] ", color.New(color.Bold).Sprintf("%s", includeOutputName)))
	utils.Print(fmt.Sprintf("%s -> %s", includeOutputName, outputFilePath), true)
	utils.PrintNewLine()
}

func printEncryptionStart(encMethod string, key string, shellcode []byte) {
	utils.PrintWhite(fmt.Sprintf("[ %s ] ", color.New(color.Bold).Sprintf("Shellcode Encryption")))
	utils.Print(fmt.Sprintf("Using: %s", color.New(color.Bold).Sprintf(encMethod)), true)
	utils.Print(fmt.Sprintf("Key: %s", color.New(color.Bold).Sprintf(key)), true)
	utils.Print(fmt.Sprintf("Size: %s bytes", color.New(color.Bold).Sprintf("%d", len(shellcode))), true)
	printShellcodeBytes("Before", shellcode)
}

func printEncryptionEnd(shellcode []byte) {
	printShellcodeBytes("After", shellcode)
	utils.PrintNewLine()
}

func printShellcodeBytes(prefix string, shellcode []byte) {
	byteArray := make([]string, len(shellcode))
	for i, b := range shellcode {
		byteArray[i] = fmt.Sprintf("0x%02X", b)
	}
	ba := fmt.Sprintf("{ %s }", strings.Join(byteArray, ", "))
	utils.Print(fmt.Sprintf("%s: %s ... }", prefix, color.New(color.Bold).Sprintf(ba[:24])), true)
}
