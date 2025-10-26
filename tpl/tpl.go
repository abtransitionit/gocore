package tpl

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// description: load txt file. and return a string
//
// Notes:
// - used to load yaml file of json file into a string
func LoadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read config template file: %v", err)
	}
	return string(data), nil
}

// // description: replace placeholders of a string from a struct
// //
// // Notes:
// // - used to load yaml file of json file into a string
func ResolveTplConfig[T any](tplFile string, vars T) (string, error) {
	tpl, err := template.New("repo").Parse(tplFile)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, vars); err != nil {
		return "", err
	}

	return buf.String(), nil
}
