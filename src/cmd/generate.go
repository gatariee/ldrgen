package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	gen "ldrgen/cmd/generator"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a loader",
	Long:  `Generate a loader based on the specified template & arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		templatePath, _ := cmd.Flags().GetString("template")
		config, err := gen.ReadConfig(filepath.Join(templatePath, "config.yaml"))
		if err != nil {
			fmt.Println("[!] there was an error reading your config, this is what we are trying to read: ", filepath.Join(templatePath, "config.yaml"))
			return
		}

		binPath, _ := cmd.Flags().GetString("bin")
		AbsBinPath, err := filepath.Abs(binPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		switch binPath {
		case "":
			cmd.Help()
			fmt.Println("\n--bin -> you must specify a .bin file to load")
			return

		default:
			_, err := os.Stat(AbsBinPath)
			if os.IsNotExist(err) {
				fmt.Println("--bin -> file does not exist, we are looking for:", AbsBinPath)
				return
			}
		}

		outputPath, _ := cmd.Flags().GetString("output")
		switch outputPath {
		case "":
			cmd.Help()
			fmt.Println("\n--output -> you must specify an output directory")
			return

		default:
			err := os.MkdirAll(outputPath, 0o755)
			if err != nil {
				fmt.Println("\n--output -> there was an error creating the output directory, please check your permissions and try again")
				return
			}
		}

		ldrToken, _ := cmd.Flags().GetString("loader")

		switch ldrToken {
		case "":
			cmd.Help()
			fmt.Println("\n--loader -> you must specify a loader token")
			return
		default:
			ldrFound := gen.TokenExists(ldrToken, config)
			if err != nil {
				fmt.Println(err)
				return
			}
			if !ldrFound {
				fmt.Println("\n--loader -> loader token not found, please check your spelling and try again")
				return
			}
		}

		enc, _ := cmd.Flags().GetString("enc")

		switch enc {
		case "":
			/* This is okay, as long as the loader doesn't require encryption */
			needEnc := gen.TokenEnc(ldrToken, config)
			if needEnc {
				fmt.Println("\n--enc -> you must specify an encryption type for loader: ", ldrToken)
				return
			}
		default:
			/* Do nothing */
		}

		arg, _ := cmd.Flags().GetString("args")
		parsed_args := make(map[string]string)
		switch arg {
		case "":
			/* This is okay, as long as the loader doesn't require arguments */

			// don't currently have a function to verify this, so we'll just let the operator do whatever they want :)

		default:
			if arg != "" {
				argList := strings.Split(arg, ",")
				parsed_args, err = gen.ParseArgs(argList)
				if err != nil {
					return
				}
			}
		}

		cleanup, _ := cmd.Flags().GetBool("cleanup")

		if binPath == "" {
			cmd.Help()
			fmt.Println("\n--bin -> you must specify a .bin file to load")
			return
		}

		if outputPath == "" {
			cmd.Help()
			return
		}

		if ldrToken == "" {
			cmd.Help()
			return
		}

		needsEnc := gen.TokenEnc(ldrToken, config)
		if needsEnc && enc == "" {
			enc = "xor"
		}

		if !needsEnc && enc != "" {
			return
		}

		err = gen.ProcessShellcodeTemplate(AbsBinPath, enc, parsed_args, outputPath, config, templatePath)
		if err != nil {
			gen.CleanShellcode(outputPath, config)
			return
		}

		err = gen.ProcessLoaderTemplate(ldrToken, enc, parsed_args, outputPath, config, templatePath)
		if err != nil {
			fmt.Println(err)
			gen.CleanShellcode(outputPath, config)
			gen.CleanLoader(outputPath, config)
			return
		}

		if enc != "" && cleanup {
			err := os.Remove(binPath + ".enc")
			if err != nil {
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("bin", "b", "", "* Path to binary file to load (required)")
	generateCmd.Flags().StringP("output", "o", "", "* Path to output directory (required)")
	generateCmd.Flags().StringP("loader", "l", "", "* Loader token to use (required)")
	generateCmd.Flags().StringP("enc", "e", "", "Encryption type to use (optional, depending on ldr)")
	generateCmd.Flags().StringP("args", "a", "", "Arguments to pass to template (optional, depending on ldr)")
	generateCmd.Flags().StringP("template", "t", "../templates", "Path to template folder (optional, default: ../templates)")
	generateCmd.Flags().BoolP("cleanup", "c", false, "Cleanup temporary files (optional, default: false)")
}
