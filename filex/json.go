package filex

import (
	"encoding/json"
	"fmt"
	"os"
)

// Description: converts 	a file into memory (as []byte) into a custom struct into memory
func LoadJsonIntoStruct[T any](data []byte) (*T, error) {
	// 1 - Define the destination struct
	var dest T

	// 2 - Unmarshal the JSON into the struct
	if err := json.Unmarshal(data, &dest); err != nil {
		return nil, fmt.Errorf("unmarshalling JSON: %w", err)
	}

	// 3 - Return a pointer to the struct
	return &dest, nil
}

// Description: reads a file from disk into a custom struct into memory
//
// Notes:
//   - use a helper function.
func LoadJsonFromFile[T any](filePath string) (*T, error) {

	// 1 - Read the entire file content into memory (as []byte no struct yet)
	dataAsByte, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", filePath, err)
	}

	// 2 - Delegate decoding
	return LoadJsonIntoStruct[T](dataAsByte)
}
