package filex

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Description: loads a YAML file from the given path into a custom struct
func LoadYamlFile[T any](filePath string) (*T, error) {
	// 1 - Read the entire file content into memory (as []byte no struct yet)
	dataAsByte, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", filePath, err)
	}
	// 2 - Delegate decoding
	return LoadYamlIntoStruct[T](dataAsByte)
}

// Description: converts a YAML file into memory (as []byte) into a custom struct T (into memory)
//
// Usage Example:
//
//	confData := []byte(`apt:\n  pkg: deb\n  ext: .list`)
//	conf, err := yamlx.LoadYamlIntoStruct[Config](confData)
//	if err != nil {
//	    log.Fatal(err)
//	}
func LoadYamlIntoStruct[T any](data []byte) (*T, error) {
	// 1 - Define the destination struct
	var dest T

	// 2 - Unmarshal the YAML into the struct
	if err := yaml.Unmarshal(data, &dest); err != nil {
		return nil, fmt.Errorf("unmarshalling YAML: %w", err)
	}

	// 3 - Return a pointer to the struct
	return &dest, nil
}

// Description: reads a YAML file, applies template rendering using a custom ctx structure and converts it back into a custom type T.
func LoadTplYamlFile[T any](filePath string, ctx any) (*T, error) {
	// 1 - Read the entire file content into memory (as []byte no struct yet)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", filePath, err)
	}
	// 2 - Delegate rendering and decoding
	return LoadTplYamlFileEmbed[T](data, ctx)
}

// Description: renders a YAML file in memory (as []byte) using a ctx structure and converts it back into a custom type T
func LoadTplYamlFileEmbed[T any](data []byte, ctx any) (*T, error) {
	// 1 - parse and execute the template
	tmpl, err := template.New("yaml").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	// 2 - execute the template
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, ctx); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}
	// 3 - delegate decoding
	return LoadYamlIntoStruct[T](rendered.Bytes())
}
