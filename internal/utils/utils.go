package utils

import (
	"os"
	"path"
	"strings"
)

func GetSecret(name string) (string, error) {
	s, err := os.ReadFile(path.Join("/run/secrets/", name))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(s)), nil
}
