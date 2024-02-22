package generate

import (
	"fmt"
)

func CountSubs(files []FileConfig) int {
	num_subs := 0
	for _, f := range files {
		num_subs += len(f.Substitutions)
	}
	return num_subs
}

func VerifySubstitutions(template_args map[string]string, files []FileConfig) error {
	for _, f := range files {
		for k := range f.Substitutions {
			if _, ok := template_args[k]; !ok {
				return fmt.Errorf("missing argument: %s", k)
			}
		}
	}
	return nil
}
