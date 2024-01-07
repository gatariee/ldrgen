package cmd 

import (
	"fmt"
	"path/filepath"
	"strings"
)

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