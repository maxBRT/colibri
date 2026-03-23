package utils

import (
	"os"
	"strings"
)

func GetSecret(name string) (string, error) {
	s, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(s)), nil
}
