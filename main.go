package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

var (
	binPath       = flag.String("bin", "", "Path to shellcode .bin file")
	outputPath    = flag.String("out", "", "Output folder")
	ldrToken      = flag.String("ldr", "", "Loader token")
	enc           = flag.String("enc", "", "Encryption type (optional)")
	args          = flag.String("args", "", "Arguments to be passed to the loader (optional)")
	template_path = flag.String("template", "", "Path to template folder (default: ./template)")
	cleanup       = flag.Bool("cleanup", false, "Cleanup? (delete encrypted shellcode file)")
	help          = flag.Bool("help", false, "Print help")
)

var Templates Config

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

  -args <args>     Arguments to be passed to the loader's template. (Optional)
                     Example: -args "key=mySecretKey1234,pid=1234"
					 
  -template <path> Path to the template folder to be used for generating the loader. (Optional)
                     Example: -template /path/to/custom/template
                     Default: ./template

  -cleanup         Cleanup? (delete encrypted shellcode file)

  -help            Print this help message.

Examples:
  ./ldr -bin ./dev/calc_shellcode/calc.bin -out ./output -ldr Inline
  ./ldr -bin ./dev/calc_shellcode/calc.bin -out ./output -ldr CreateRemoteThread 
  ./ldr -bin ./dev/calc_shellcode/calc.bin -out ./output -ldr CreateThread

  ./ldr -bin ./dev/calc_shellcode/calc.bin -out ./output -ldr Inline_Xor -enc xor -args "key=test" -cleanup
  ./ldr -bin ./dev/calc_shellcode/calc.bin -out ./output -ldr CreateThread_Xor -enc xor -args "key=test" -cleanup  
  ./ldr -bin ./dev/calc_shellcode/calc.bin -out ./output -ldr CreateThread_Xor_Sleep -enc xor -args "key=test, sleep=5" -cleanup
`
	fmt.Println(text)
}

func tokenMethod(token string) string {
	for _, t := range Templates.Loader {
		if strings.EqualFold(t.Token, token) {
			return t.Method
		}
	}
	return ""
}

func tokenExists(token string) bool {
	for _, t := range Templates.Loader {
		if strings.EqualFold(t.Token, token) {
			return true
		}
	}
	return false
}

func tokenEnc(token string) bool {
	for _, t := range Templates.Loader {
		if strings.EqualFold(t.Token, token) {
			return t.KeyRequired
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
		// TODO: this can error out if key is not provided, remember to validate

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
		/* I'm sorry */

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

func parseArgs(arg []string) (map[string]string, error) {
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

func countSubs(files []FileConfig) int {
	num_subs := 0
	for _, f := range files {
		num_subs += len(f.Substitutions)
	}
	return num_subs
}

func verifySubstitutions(template_args map[string]string, files []FileConfig) error {
	for _, f := range files {
		for k := range f.Substitutions {
			if _, ok := template_args[k]; !ok {
				return fmt.Errorf("missing argument: %s", k)
			}
		}
	}
	return nil
}

func ProcessLoaderTemplate(token string, enc string, template_args map[string]string) error {
	method := tokenMethod(token)
	fmt.Println("[LDR] Using:", method)

	for _, l := range Templates.Loader {
		if strings.EqualFold(l.Token, token) {

			// check that number of substitutions, match the number of arguments

			if len(template_args) != countSubs(l.Files) {
				return fmt.Errorf("number of arguments does not match number of substitutions: expected %d, got %d", countSubs(l.Files), len(template_args))
			}

			// check that all substitutions are present in the arguments
			err := verifySubstitutions(template_args, l.Files)
			if err != nil {
				return err
			}

			for _, f := range l.Files {
				content, err := ReadFile(filepath.Join(*template_path, f.SourcePath))
				if err != nil {
					return err
				}

				for k, v := range template_args {
					for s, r := range f.Substitutions {
						if strings.EqualFold(k, s) {
							fmt.Println("[TEMPLATING] Replacing:", r, "with:", v)
							content = strings.ReplaceAll(content, r, v)
						}
					}
				}

				err = SaveToFile(*outputPath, f.OutputPath, content)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}

	return nil
}

func readConfig(path string) (*Config, error) {
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

func init() {
	flag.Usage = PrintHelp
	flag.Parse()

	if *help {
		PrintHelp()
		os.Exit(0)
	}

	if *binPath == "" || *outputPath == "" || *ldrToken == "" {
		PrintHelp()
		os.Exit(0)
	}

	if *template_path == "" {
		fmt.Println("[*] Using default template path: ./templates")
		*template_path = "./templates"
	}

	config, err := readConfig(filepath.Join(*template_path, "config.yaml"))
	if err != nil {
		fmt.Println("[!] Error reading config file:", err)
		os.Exit(1)
	}

	Templates = *config

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
}

func cleanShellcode() {
	fmt.Println("deleting -", *outputPath+"/"+Templates.Shellcode.SourceOutputName)
	err := os.Remove(*outputPath + "/" + Templates.Shellcode.SourceOutputName)
	if err != nil {
		return
	}

	fmt.Println("deleting -", *outputPath+"/"+Templates.Shellcode.IncludeOutputName)
	err = os.Remove(*outputPath + "/" + Templates.Shellcode.IncludeOutputName)
	if err != nil {
		return
	}
}

func cleanLoader() {
	for _, l := range Templates.Loader {
		for _, f := range l.Files {
			fmt.Println("deleting -", *outputPath+"/"+f.OutputPath)
			err := os.Remove(*outputPath + "/" + f.OutputPath)
			if err != nil {
				return
			}
		}
	}
}

func main() {
	template_args := make(map[string]string)
	if *args != "" {
		var err error
		template_args, err = parseArgs(strings.Split(*args, ","))
		if err != nil {
			fmt.Println("[!] Error parsing arguments:", err)
			os.Exit(1)
		}
	}

	err := ProcessShellcodeTemplate(*binPath, *enc, template_args)
	if err != nil {
		fmt.Println("Error processing shellcode:", err)

		cleanShellcode()

		os.Exit(1)
	}

	err = ProcessLoaderTemplate(strings.ToLower(*ldrToken), *enc, template_args)
	if err != nil {
		fmt.Println("Error processing loader:", err)

		cleanShellcode()
		cleanLoader()

		os.Exit(1)
	}

	if *enc != "" && *cleanup {
		fmt.Println("[CLEANUP] Removing encrypted shellcode file:", *binPath+".enc")
		err = os.Remove(*binPath + ".enc")
		if err != nil {
			fmt.Println("[!] Error removing encrypted shellcode file:", err)
			return
		}
	}

	fmt.Println("[*] Done!")
}
