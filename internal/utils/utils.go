package utils

import (
	"io"
	"os"
	"strconv"
	"strings"
)

func GetSecret(name string) (string, error) {
	rootPath := "/run/secrets"
	r, err := os.OpenRoot(rootPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	name = strings.TrimPrefix(name, rootPath)
	name = strings.TrimPrefix(name, "/")

	f, err := r.Open(name)
	if err != nil {
		return "", err
	}

	s, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(s)), nil
}

func ParseBool(val string) bool {
	b, err := strconv.ParseBool(val)
	return err == nil && b
}
