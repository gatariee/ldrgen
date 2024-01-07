package cmd

import (
	"fmt"
)

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
  ./ldr -bin ./dev/calc_shellcode/calc.bin -out ./output -ldr EarlyBirdAPC -args "pname=c:\\\windows\\\system32\\\cmd.exe"
  ./ldr -bin ./dev/calc_shellcode/calc.bin -out ./output -ldr EarlyBirdAPC_Buffed -args "pname=c:\\windows\\system32\\cmd.exe, key=mySecretKey1234"
`
	fmt.Println(text)
}
