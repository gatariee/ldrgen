package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

)

var (
	binPath       = flag.String("bin", "", "Path to shellcode .bin file")
	outputPath    = flag.String("out", "", "Output folder")
	ldrToken      = flag.String("ldr", "", "Loader token")
	enc           = flag.String("enc", "", "Encryption type (optional)")
	args          = flag.String("args", "", "Arguments to be passed to the loader (optional)")
	template_path = flag.String("template", "", "Path to template folder (default: ./template)")
	cleanup       = flag.Bool("cleanup", false, "Cleanup? (delete encrypted shellcode file)")
	help          = flag.Bool("help", false, "Print help")
	Templates     Config
)

func parseArgs(arg []string) (map[string]string, error) {
	/*
	Handles the parsing of specifically the "args" flag, which is a comma-separated list of key=value pairs.
	Example: ... -args "key1=value1,key2=value2,key3=value3"

	*/
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

func init() {

	/* 
	I hate the way this handles incorrect usage, but I'm way too lazy to improve it. So, we're handling invalid cases manually here.
	*/

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

func Execute() {
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
