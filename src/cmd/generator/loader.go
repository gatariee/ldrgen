package generate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"ldrgen/cmd/utils"
)

func ProcessLoaderTemplate(token string, enc string, template_args map[string]string, outputPath string, config *Config, template_path string) error {
	method := TokenMethod(token, config)

	message := fmt.Sprintf("[ %s ] ", color.New(color.Bold).Sprintf("Loader"))
	utils.PrintWhite(message)

	message = fmt.Sprintf("Using: %s", color.New(color.Bold).Sprintf(token))
	utils.Print(message, true)

	message = fmt.Sprintf("Calls: %s", color.New(color.Bold).Sprintf(method))
	utils.Print(message, true)

	for _, l := range config.Loader {
		if strings.EqualFold(l.Token, token) {

			// check that number of substitutions, match the number of arguments
			if len(template_args) != CountSubs(l.Files) {
				return fmt.Errorf("number of arguments does not match number of substitutions: expected %d, got %d", CountSubs(l.Files), len(template_args))
			}

			// check that all substitutions are present in the arguments
			err := VerifySubstitutions(template_args, l.Files)
			if err != nil {
				return err
			}

			utils.PrintNewLine()
			message = fmt.Sprintf("[ %s ] ", color.New(color.Bold).Sprintf("Substitution"))
			utils.PrintWhite(message)

			for _, f := range l.Files {
				content, err := ReadFile(filepath.Join(template_path, f.SourcePath))
				if err != nil {
					return err
				}

				for k, v := range template_args {
					for s, r := range f.Substitutions {
						if strings.EqualFold(k, s) {
							message = fmt.Sprintf("%s -> %s", color.New(color.Bold).Sprintf(f.SourcePath), color.New(color.Bold).Sprintf(f.OutputPath))
							utils.Print(message, true)
							message = fmt.Sprintf("\t%s -> %s", r, v)
							utils.Print(message, true)
							content = strings.ReplaceAll(content, r, v)
							utils.PrintNewLine()
						}
					}
				}

				err = SaveToFile(outputPath, f.OutputPath, content)
				if err != nil {
					return err
				}

			}
			return nil
		}
	}

	/* This should *never* happen */
	return fmt.Errorf("loader not found")
}
