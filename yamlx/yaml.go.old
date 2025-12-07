package yamlx

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Description: loads a YAML file from the given path into a struct of type T.
//
// Usage Example:
//
//	conf, err := yamlx.LoadYamlFile[Config]("conf.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
func LoadYamlFile[T any](filePath string) (*T, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", filePath, err)
	}
	return LoadYamlFileEmbed[T](data)
}

// Description: unmarshals YAML from a byte slice into a struct of type T.
//
// Usage Example:
//
//	confData := []byte(`apt:\n  pkg: deb\n  ext: .list`)
//	conf, err := yamlx.LoadYamlFileEmbed[Config](confData)
//	if err != nil {
//	    log.Fatal(err)
//	}
func LoadYamlFileEmbed[T any](data []byte) (*T, error) {
	var cfg T
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling YAML: %w", err)
	}
	return &cfg, nil
}

// Description: reads a YAML file, applies template rendering with ctx, and unmarshals into T.
func LoadTplYamlFile[T any](filePath string, ctx any) (*T, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", filePath, err)
	}
	return LoadTplYamlFileEmbed[T](data, ctx)
}

// Description: renders embedded YAML (as []byte) with ctx and unmarshals into T.
func LoadTplYamlFileEmbed[T any](data []byte, ctx any) (*T, error) {
	tmpl, err := template.New("yaml").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, ctx); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return LoadYamlFileEmbed[T](rendered.Bytes())
}
