package cmd 

import (
	"strings"
)

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