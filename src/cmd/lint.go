package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"ldrgen/cmd/generator"
	"ldrgen/cmd/utils"

	"github.com/chzyer/readline"

	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint a template for errors and warnings",
	Long:  `Lint a template for errors and warnings`,
	Run: func(cmd *cobra.Command, args []string) {
		template_path, _ := cmd.Flags().GetString("template")

		utils.PrintBanner()

		message := fmt.Sprintf("Linting template: '%s'", color.New(color.Bold).Sprintf("%s", template_path))
		utils.PrintInfo(message, true)

		/*
			@ generate.ReadConfig(template_path) can error out, so we'll rewrite a function with error handling and linting in mind
		*/

		cfg, err := generate.ReadConfig(template_path)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Error reading config: %s", err), true)
			return
		}

		num_ldrs := len(cfg.Loader)
		if num_ldrs == 0 {
			utils.PrintError("No loaders found in the config", true)
			return
		}

		message = fmt.Sprintf("Successfully read %s loader(s) from the config", color.New(color.Bold).Sprintf("%d", num_ldrs))
		utils.PrintSuccess(message, true)

		/* Ldr path is always in the same cwd as the config */
		ldr_path := strings.TrimSuffix(template_path, "config.yaml")

		utils.PrintNewLine()

		rl, err := readline.New("> ")
		if err != nil {
			utils.PrintError(fmt.Sprintf("Error creating readline: %s", err), true)
			return
		}
		defer rl.Close()

		PrintAllLdrs(cfg)

		for {
			utils.PrintNewLine()

			utils.PrintInfo("Enter a token to lint", true)
			line, err := rl.Readline()
			if err != nil {
				utils.PrintError(fmt.Sprintf("Error reading input: %s", err), true)
				return
			}

			switch line {
			case "clear":
				utils.ClearScreen()
				continue

			case "help":
				utils.PrintNewLine()
				utils.PrintInfo("Available commands:", true)
				utils.Print("- clear: clear the screen", true)
				utils.Print("- help: show this help message", true)
				utils.Print("- ldrs: list all loader tokens", true)
				utils.Print("- exit: exit the linting session", true)

				continue

			case "exit":
				return

			case "ldrs":

				PrintAllLdrs(cfg)
				continue

			case "":
				continue
			}

			if !generate.TokenExists(line, cfg) {
				utils.PrintError(fmt.Sprintf("Token '%s' not found in the config", line), true)
				continue
			}

			message = fmt.Sprintf("Token information for: %s", color.New(color.Bold).Sprintf("%s", line))
			utils.PrintInfo(message, true)
			PrintLdrInformation(cfg, line)

			err = VerifySubstitutions(cfg, line, ldr_path)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Error verifying substitutions: %s", err), true)
				return
			}
		}
	},
}

func VerifySubstitutions(cfg *generate.Config, token string, ldr_path string) error {
	utils.PrintNewLine()

	for _, l := range cfg.Loader {
		if l.Token == token {
			for _, f := range l.Files {
				if len(f.Substitutions) != 0 {
					message := fmt.Sprintf("Found %s substitutions for file: '%s'", color.New(color.Bold).Sprintf("%d", len(f.Substitutions)), color.New(color.Bold).Sprintf("%s", f.SourcePath))
					utils.PrintInfo(message, true)

					content, err := os.ReadFile(ldr_path + f.SourcePath)
					if err != nil {
						return err
					}
					message = fmt.Sprintf("Successfully read %s bytes from file: '%s'", color.New(color.Bold).Sprintf("%d", len(content)), color.New(color.Bold).Sprintf("%s", f.SourcePath))
					utils.PrintSuccess(message, true)

					for _, v := range f.Substitutions {
						utils.PrintNewLine()
						message = fmt.Sprintf("Checking: '%s'", color.New(color.Bold).Sprintf("%s", v))
						utils.PrintInfo(message, true)
						if !strings.Contains(string(content), v) {
							return fmt.Errorf("substitution '%s' not found in file '%s'", v, f.SourcePath)
						}
						message = fmt.Sprintf("%s: OK", color.New(color.Bold).Sprintf("%s", v))
						utils.PrintSuccess(message, true)
					}
				}
			}
		}
	}

	return nil
}

func PrintAllLdrs(cfg *generate.Config) {
	for _, l := range cfg.Loader {
		fmt.Println("Token:", l.Token)
	}
}

func PrintLdrInformation(cfg *generate.Config, token string) {
	for _, l := range cfg.Loader {
		if l.Token == token {
			message := fmt.Sprintf("\tToken: %s", color.New(color.Bold).Sprintf("%s", l.Token))
			utils.Print(message, true)
			message = fmt.Sprintf("\tMethod: %s", color.New(color.Bold).Sprintf("%s", l.Method))
			utils.Print(message, true)
			message = fmt.Sprintf("\tKey Required: %t", l.KeyRequired)
			utils.Print(message, true)

			if l.EncType != nil {
				message = fmt.Sprintf("\tEnc Type: %s", color.New(color.Bold).Sprintf("%s", *l.EncType))
				utils.Print(message, true)
			}

			for _, f := range l.Files {
				message = fmt.Sprintf("\tFile: %s", color.New(color.Bold).Sprintf("%s", f.SourcePath))
				utils.Print(message, true)
				message = fmt.Sprintf("\t\tOutput: %s", color.New(color.Bold).Sprintf("%s", f.OutputPath))
				utils.Print(message, true)

				if len(f.Substitutions) != 0 {
					message = fmt.Sprintf("\t\tSubstitutions: %d", len(f.Substitutions))
					utils.Print(message, true)
					for k, v := range f.Substitutions {
						message = fmt.Sprintf("\t\t\t%s: %s", color.New(color.Bold).Sprintf("%s", k), color.New(color.Bold).Sprintf("%s", v))
						utils.Print(message, true)
					}
				}
			}
		}
	}
}

func init() {
	lintCmd.Flags().StringP("template", "t", "", "Path to the template to lint")
	lintCmd.MarkFlagRequired("template")
}
