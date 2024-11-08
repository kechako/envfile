// Package envfile provides functionality to parse files containing environment variables in the format key=value.
package envfile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"iter"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Envs represents an array of strings, where each string follows the format key=value.
type Envs []string

// Envs returns an iterator over key-value pairs from Envs.
func (envs Envs) Envs() iter.Seq2[string, string] {
	return func(yield func(key, value string) bool) {
		for _, line := range envs {
			key, value, _ := strings.Cut(line, "=")
			if !yield(key, value) {
				return
			}
		}
	}
}

// Map converts Envs array into a map[string]string.
func (envs Envs) Map() map[string]string {
	m := map[string]string{}

	for key, value := range envs.Envs() {
		m[key] = value
	}

	return m
}

// ParseFile reads the named file, parses the content in the key=value format, and returns an Envs type.
func ParseFile(name string) (Envs, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(file)
}

var utf8BOM = [3]byte{0xef, 0xbb, 0xbf}

// Parse reads data from the r, parses the content in the key=value format, and returns an Envs type.
func Parse(r io.Reader) (Envs, error) {
	var envs []string
	s := bufio.NewScanner(r)

	for lineNo := 1; s.Scan(); lineNo++ {
		b := s.Bytes()
		if !utf8.Valid(b) {
			return nil, fmt.Errorf("invalid UTF-8 bytes at line %d", lineNo)
		}

		// remove the UTF-8 BOM from the first line
		if lineNo == 1 {
			b = bytes.TrimPrefix(b, utf8BOM[:])
		}

		// remove all leading whitespace from each line
		line := strings.TrimLeftFunc(string(b), unicode.IsSpace)

		// ignore empty lines or comment lines start with '#'
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		key, _, _ := strings.Cut(line, "=")
		if len(key) == 0 {
			return nil, fmt.Errorf("no variable key on line %d", lineNo)
		}

		if strings.ContainsFunc(key, unicode.IsSpace) {
			return nil, fmt.Errorf("the variable key contains whitespace: '%s'", key)
		}

		envs = append(envs, line)
	}

	return envs, nil
}

// LoadEnvs sets environment variables from envs.
func LoadEnvs(envs Envs) error {
	for key, value := range envs.Envs() {
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// Parse reads data from the r, parses the content in the key=value format, and sets environment variables.
func Load(r io.Reader) error {
	envs, err := Parse(r)
	if err != nil {
		return err
	}
	err = LoadEnvs(envs)
	if err != nil {
		return err
	}
	return nil
}

// ParseFile reads the named file, parses the content in the key=value format, and sets environment variables.
func LoadFile(name string) error {
	envs, err := ParseFile(name)
	if err != nil {
		return err
	}
	err = LoadEnvs(envs)
	if err != nil {
		return err
	}
	return nil
}
