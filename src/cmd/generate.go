package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	gen "ldrgen/cmd/generator"
	utils "ldrgen/cmd/utils"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a loader",
	Long:  `Generate a loader based on the specified template & arguments`,
	Run: func(cmd *cobra.Command, args []string) {
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

		templatePath, _ := cmd.Flags().GetString("template")
		config, err := gen.ReadConfig(filepath.Join(templatePath, "config.yaml"))
		if err != nil {
			fmt.Println("[!] there was an error reading your config, this is what we are trying to read: ", filepath.Join(templatePath, "config.yaml"))
			return
		}
		ldrFound := gen.TokenExists(ldrToken, config)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !ldrFound {
			fmt.Println("\n--loader -> loader token not found, please check your spelling and try again")
			return
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

		absOutputPath, err := filepath.Abs(outputPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		utils.PrintSuccess("Looks like everything worked!", true)
		message := fmt.Sprintf("Your loader is located at: %s", color.New(color.Bold).Sprintf(absOutputPath))
		utils.PrintInfo(message, true)

		utils.PrintNewLine()

		compile, _ := cmd.Flags().GetString("compile")
		err = gen.CompileLoader(absOutputPath, compile)
		if err != nil {
			fmt.Println(err)
			return
		}
		utils.PrintSuccess("Loader compiled successfully!", true)
	},
}

func init() {
	generateCmd.Flags().StringP("bin", "b", "", "* Path to binary file to load")
	generateCmd.MarkFlagRequired("bin")
	generateCmd.Flags().StringP("output", "o", "", "* Path to output directory")
	generateCmd.MarkFlagRequired("output")
	generateCmd.Flags().StringP("loader", "l", "", "* Loader token to use")
	generateCmd.MarkFlagRequired("loader")
	generateCmd.Flags().StringP("template", "t", "", "Path to template directory")
	generateCmd.MarkFlagRequired("template")

	generateCmd.Flags().StringP("enc", "e", "", "Encryption type to use")
	generateCmd.Flags().StringP("args", "a", "", "Arguments to pass to template")

	generateCmd.Flags().BoolP("cleanup", "c", false, "Cleanup temporary files")
	generateCmd.Flags().StringP("compile", "C", "", "Compile loader with x86 or x64")
}
