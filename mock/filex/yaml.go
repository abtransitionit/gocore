package filex

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

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
