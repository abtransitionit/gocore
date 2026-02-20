package filex

import (
	"fmt"
	"os"

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

// description: loads a YAML file into memory (as []byte) into a custom struct T
//
// Notes:
// - the YAML file can be external (provided by the user)
// - the YAML file fallback to an embedded (auto cached)
func LoadYamlIntoStruct2[T any](embedded []byte, externalPath string) (*T, error) {
	var data []byte

	// 1 - the external file is provided
	if externalPath != "" {
		exists, err := ExistsFile(externalPath)
		if err != nil {
			// 2 - the external cannot be accessed
			return nil, fmt.Errorf("getting external YAML %q: %w", externalPath, err)
		}
		// 3 - the external file exists → load it
		if exists {
			b, err := os.ReadFile(externalPath)
			if err != nil {
				// 4 - the external file cannot be read
				return nil, fmt.Errorf("cannot read external YAML %q: %w", externalPath, err)
			}
			data = b
		}
	}

	// 5 - the external file is not provided or does not exist → fallback to embedded
	if data == nil {
		data = embedded
	}

	// 6 - Unmarshal YAML into struct (embedded or external)
	cfg, err := LoadYamlIntoStruct[T](data)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling YAML: %w", err)
	}

	return cfg, nil
}

// description: loads a YAML file into memory (as []byte) into a custom struct T
//
// Notes:
// - the YAML file is external (provided by the user)
func LoadExternalYamlIntoStruct[T any](externalPath string) (*T, error) {
	var data []byte

	// 1 - Check if external path is provided
	if externalPath == "" {
		return nil, fmt.Errorf("external YAML path not provided")
	}

	// 2 - Check if external file exists and is accessible
	exists, err := ExistsFile(externalPath)
	if err != nil {
		return nil, fmt.Errorf("getting external YAML %q: %w", externalPath, err)
	}

	// 3 - Fail if the file does not exist
	if !exists {
		return nil, fmt.Errorf("external YAML %q does not exist", externalPath)
	}

	// 4 - Read the external file
	b, err := os.ReadFile(externalPath)
	if err != nil {
		return nil, fmt.Errorf("reading existing external YAML %q: %w", externalPath, err)
	}
	data = b

	// 5 - Unmarshal YAML into struct
	cfg, err := LoadYamlIntoStruct[T](data)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling YAML %q: %w", externalPath, err)
	}

	return cfg, nil
}

func LoadYamlAsByte2(embedded []byte, externalPath string) ([]byte, error) {
	var data []byte

	// 1 - the external file is provided
	if externalPath != "" {
		exists, err := ExistsFile(externalPath)
		if err != nil {
			// 2 - the external cannot be accessed
			return nil, fmt.Errorf("getting external YAML %q: %w", externalPath, err)
		}

		// 3 - the external file exists → load it
		if exists {
			b, err := os.ReadFile(externalPath)
			if err != nil {
				// 4 - the external file cannot be read
				return nil, fmt.Errorf("cannot read external YAML %q: %w", externalPath, err)
			}
			data = b
		}
	}

	// 5 - the external file is not provided or does not exist → fallback to embedded
	if data == nil {
		data = embedded
	}

	return data, nil
}
func LoadExternalYamlAsByte(externalPath string) ([]byte, error) {
	// 1 - Check if external path is provided
	if externalPath == "" {
		return nil, fmt.Errorf("external YAML path not provided")
	}

	// 2 - Check if external file exists and is accessible
	exists, err := ExistsFile(externalPath)
	if err != nil {
		return nil, fmt.Errorf("getting external YAML %q: %w", externalPath, err)
	}

	// 3 - Fail if the file does not exist
	if !exists {
		return nil, fmt.Errorf("external YAML %q does not exist", externalPath)
	}

	// 4 - Read the external file
	data, err := os.ReadFile(externalPath)
	if err != nil {
		return nil, fmt.Errorf("reading existing external YAML %q: %w", externalPath, err)
	}

	return data, nil
}
