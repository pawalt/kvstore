package kv

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// alphanumeric separated by /
	r, _ = regexp.Compile(`^[a-z0-9]+(\/[a-z0-9]+)*$`)
)

func ValidPath(path string) bool {
	return r.Match([]byte(path))
}

func ParsePath(path string) ([]string, error) {
	if !ValidPath(path) {
		return nil, errors.New("path " + path + " is invalid")
	}

	return strings.Split(path, "/"), nil
}
