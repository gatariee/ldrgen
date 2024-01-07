package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ProcessShellcodeTemplate(binPath string, enc string, args map[string]string) error {
	if enc != "" {
		_, ok := args["key"]
		if !ok {
			return fmt.Errorf("key not provided")
		}
	}

	switch enc {
	case "xor":
		key := args["key"]
		fmt.Println("[!!!] XORing shellcode with key:", key)
		shellcode, err := os.ReadFile(binPath)
		if err != nil {
			return err
		}

		var byteArray []string
		for _, b := range shellcode {
			byteArray = append(byteArray, fmt.Sprintf("0x%02X", b))
		}

		ba := fmt.Sprintf("{ %s }", strings.Join(byteArray, ", "))

		fmt.Println("[Sanity Check] Starting bytes (before XOR): ", ba[0:24], "... }")

		for i := 0; i < len(shellcode); i++ {
			shellcode[i] ^= key[i%len(key)]
		}

		err = os.WriteFile(binPath+".enc", shellcode, 0o644)
		if err != nil {
			return err
		}

		binPath = binPath + ".enc"

		fmt.Printf("[*] Encrypted shellcode saved to: %s\n\n", binPath)

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

	shellcode_template, err := ReadFile(filepath.Join(*template_path, Templates.Shellcode.SourcePath))
	if err != nil {
		return err
	}

	shellcodeTemplate := strings.ReplaceAll(shellcode_template, Templates.Shellcode.Substitutions["shellcode"], shellcodeArray)
	shellcodeTemplate = strings.ReplaceAll(shellcodeTemplate, Templates.Shellcode.Substitutions["shellcode_size"], fmt.Sprintf("%d", fileInfo.Size()))

	abs, err := filepath.Abs(strings.TrimSpace(*outputPath))
	if err != nil {
		return err
	}

	fmt.Println("[*] Saving generated files to:", abs)

	fmt.Println("[SHELLCODE] Shellcode size:", fileInfo.Size(), "bytes")
	fmt.Println("[SHELLCODE] Starting bytes: ", shellcodeArray[0:24], "... }")

	err = SaveToFile(*outputPath, Templates.Shellcode.SourceOutputName, shellcodeTemplate)
	if err != nil {
		return err
	}

	fmt.Println("[*] Shellcode.c -> OK")

	header_template, err := ReadFile(filepath.Join(*template_path, Templates.Shellcode.IncludePath))
	if err != nil {
		return err
	}
	fmt.Println("[*] Shellcode.h -> OK")

	err = SaveToFile(*outputPath, Templates.Shellcode.IncludeOutputName, header_template)
	if err != nil {
		return err
	}

	return nil
}
