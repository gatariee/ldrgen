package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	gen "ldrgen/cmd/generator"
	utils "ldrgen/cmd/utils"

	"github.com/spf13/cobra"
)

type Config struct {
	Name        string   `json:"name"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Template    Template `json:"template"`
	Arch        string   `json:"arch"`
	Compile     Compile  `json:"compile"`
	OutputDir   string   `json:"output_dir"`
}

type Template struct {
	Path          string            `json:"path"`
	Token         string            `json:"token"`
	EncType       string            `json:"enc_type"`
	Substitutions map[string]string `json:"substitutions"`
}

type Compile struct {
	Automatic bool              `json:"automatic"`
	Make      string            `json:"make"`
	Gcc       map[string]string `json:"gcc"`
	Strip     map[string]string `json:"strip"`
}

func OutputDirExists(outputDir string) bool {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return false
	}

	return true
}

func GetAbs(filePath string) (string, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	return absFilePath, nil
}

func PatchMakefile(x64 string, x86 string, makefilePath string) error {
	x64_default := "x86_64-w64-mingw32-gcc"
	x86_default := "i686-w64-mingw32-gcc"

	makefile, err := os.ReadFile(makefilePath)
	if err != nil {
		return err
	}

	makefileString := string(makefile)
	makefileString = strings.Replace(makefileString, x64_default, x64, -1)
	makefileString = strings.Replace(makefileString, x86_default, x86, -1)

	err = os.WriteFile(makefilePath, []byte(makefileString), 0o644)
	if err != nil {
		return err
	}

	return nil
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Generate a loader from a profile",
	Long:  `Generate a loader from a profile`,
	Run: func(cmd *cobra.Command, args []string) {
		profilePath, _ := cmd.Flags().GetString("profile")
		shellcodePath, _ := cmd.Flags().GetString("shellcode")

		utils.PrintInfo("Loading profile: "+profilePath, true)

		profileFile, err := os.ReadFile(profilePath)
		if err != nil {
			log.Fatal(err)
		}

		var config Config
		err = json.Unmarshal(profileFile, &config)
		if err != nil {
			log.Fatal(err)
		}

		utils.PrintNewLine()

		utils.PrintSuccess(fmt.Sprintf(
			"%s - %s",
			config.Name,
			config.Author,
		), true)

		AbsShellcodePath, err := GetAbs(shellcodePath)
		if err != nil {
			log.Fatal(err)
		}

		template, err := gen.ReadConfig(config.Template.Path)
		if err != nil {
			log.Fatal(err)
		}

		BaseTemplatePath := filepath.Dir(config.Template.Path)
		if BaseTemplatePath == "" {
			BaseTemplatePath = "."
		}

		if !OutputDirExists(config.OutputDir) {
			err := os.MkdirAll(config.OutputDir, 0o755)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = gen.ProcessShellcodeTemplate(
			AbsShellcodePath,
			config.Template.EncType,
			config.Template.Substitutions,
			config.OutputDir,
			template,
			BaseTemplatePath,
		)
		if err != nil {
			log.Fatal(err)
		}

		err = gen.ProcessLoaderTemplate(
			config.Template.Token,
			config.Template.EncType,
			config.Template.Substitutions,
			config.OutputDir,
			template,
			BaseTemplatePath,
		)

		if err != nil {
			log.Fatal(err)
		}

		err = PatchMakefile(
			config.Compile.Gcc["x64"],
			config.Compile.Gcc["x86"],
			filepath.Join(config.OutputDir, "makefile"),
		)

		if err != nil {
			log.Fatal(err)
		}

		absSrcPath, err := GetAbs(config.OutputDir)
		if err != nil {
			log.Fatal(err)
		}

		utils.PrintSuccess("Looks like everything worked!", true)
		message := fmt.Sprintf("Source files located at: %s", absSrcPath)
		utils.PrintInfo(message, true)

		if config.Compile.Automatic {
			utils.PrintNewLine()
			utils.PrintInfo("Compiling loader...", true)
			err = gen.CompileLoader(config.OutputDir, config.Arch, config.Compile.Make)
			if err != nil {
				log.Fatal(err)
			}
			implantPath := filepath.Join(absSrcPath, "bin", fmt.Sprintf("implant_%s.exe", config.Arch))
			utils.PrintSuccess(fmt.Sprintf("-> %s", implantPath), true)
		}
	},
}

func init() {
	profileCmd.Flags().StringP("profile", "p", "", "Path to profile (e.g /home/user/ldrgen/profiles/crt.json)")
	profileCmd.MarkFlagRequired("profile")

	profileCmd.Flags().StringP("shellcode", "s", "", "Path to shellcode (e.g /home/user/cobalt/payloads/beacon_x64.bin)")
	profileCmd.MarkFlagRequired("shellcode")
}
