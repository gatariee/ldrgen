package generate

import (
	"strings"
)

func TokenMethod(token string, config *Config) string {
	for _, t := range config.Loader {
		if strings.EqualFold(t.Token, token) {
			return t.Method
		}
	}
	return ""
}

func TokenExists(token string, config *Config) bool {
	for _, t := range config.Loader {
		if strings.EqualFold(t.Token, token) {
			return true
		}
	}
	return false
}

func TokenEnc(token string, config *Config) bool {
	for _, t := range config.Loader {
		if strings.EqualFold(t.Token, token) {
			return t.KeyRequired
		}
	}
	return false
}
