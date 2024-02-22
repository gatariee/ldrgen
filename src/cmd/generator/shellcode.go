package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"ldrgen/cmd/utils"
)

func ProcessShellcodeTemplate(binPath string, enc string, args map[string]string, outputPath string, config *Config, template_path string) error {
	if enc != "" {
		_, ok := args["key"]
		if !ok {
			return fmt.Errorf("key not provided")
		}
	}

	switch enc {
	case "xor":
		key := args["key"]

		message := fmt.Sprintf("[ %s ] ", color.New(color.Bold).Sprintf("Shellcode Encryption"))
		utils.PrintWhite(message)

		message = fmt.Sprintf("Using: %s", color.New(color.Bold).Sprintf("XOR"))
		utils.Print(message, true)

		message = fmt.Sprintf("Key: %s", color.New(color.Bold).Sprintf(key))
		utils.Print(message, true)

		shellcode, err := os.ReadFile(binPath)
		if err != nil {
			return err
		}

		var byteArray []string
		for _, b := range shellcode {
			byteArray = append(byteArray, fmt.Sprintf("0x%02X", b))
		}

		ba := fmt.Sprintf("{ %s }", strings.Join(byteArray, ", "))

		message = fmt.Sprintf("Size: %s bytes", color.New(color.Bold).Sprintf("%d", len(shellcode)))
		utils.Print(message, true)

		message = fmt.Sprintf("Before: %s ... }", color.New(color.Bold).Sprintf(ba[0:24]))
		utils.Print(message, true)

		for i := 0; i < len(shellcode); i++ {
			shellcode[i] ^= key[i%len(key)]
		}

		byteArray = nil
		for _, b := range shellcode {
			byteArray = append(byteArray, fmt.Sprintf("0x%02X", b))
		}

		ba = fmt.Sprintf("{ %s }", strings.Join(byteArray, ", "))
		message = fmt.Sprintf("After: %s ... }", color.New(color.Bold).Sprintf(ba[0:24]))
		utils.Print(message, true)

		err = os.WriteFile(binPath+".enc", shellcode, 0o644)
		if err != nil {
			return err
		}

		binPath = binPath + ".enc"

		utils.PrintNewLine()
	default:
	}

	shellcodeArray, err := ToCArray(binPath)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(binPath)
	if err != nil {
		return err
	}

	shellcode_template, err := ReadFile(filepath.Join(template_path, config.Shellcode.SourcePath))
	if err != nil {
		return err
	}

	shellcodeTemplate := strings.ReplaceAll(shellcode_template, config.Shellcode.Substitutions["shellcode"], shellcodeArray)
	shellcodeTemplate = strings.ReplaceAll(shellcodeTemplate, config.Shellcode.Substitutions["shellcode_size"], fmt.Sprintf("%d", fileInfo.Size()))

	abs, err := filepath.Abs(strings.TrimSpace(outputPath))
	if err != nil {
		return err
	}

	message := fmt.Sprintf("[ %s ] ", color.New(color.Bold).Sprintf("%s", config.Shellcode.SourceOutputName))
	utils.PrintWhite(message)

	message = fmt.Sprintf("Size: %s bytes", color.New(color.Bold).Sprintf("%d", fileInfo.Size()))
	utils.Print(message, true)

	message = fmt.Sprintf("Starting bytes: %s ... }", color.New(color.Bold).Sprintf("%s", shellcodeArray[0:24]))
	utils.Print(message, true)

	err = SaveToFile(outputPath, config.Shellcode.SourceOutputName, shellcodeTemplate)
	if err != nil {
		return err
	}

	message = fmt.Sprintf("%s -> %s", config.Shellcode.SourceOutputName, abs+"/"+config.Shellcode.SourceOutputName)
	utils.Print(message, true)

	header_template, err := ReadFile(filepath.Join(template_path, config.Shellcode.IncludePath))
	if err != nil {
		return err
	}

	message = fmt.Sprintf("[ %s ] ", color.New(color.Bold).Sprintf("%s", config.Shellcode.IncludeOutputName))
	utils.PrintNewLine()
	utils.PrintWhite(message)

	err = SaveToFile(outputPath, config.Shellcode.IncludeOutputName, header_template)
	if err != nil {
		return err
	}

	message = fmt.Sprintf("%s -> %s", config.Shellcode.IncludeOutputName, abs+"/"+config.Shellcode.IncludeOutputName)
	utils.Print(message, true)

	utils.PrintNewLine()

	return nil
}
