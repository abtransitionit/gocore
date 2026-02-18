package filex

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// description: build a absolute file path from a relative path
//
// note:
// - default base path is the user home directory
func GetUserFilePath(relPath string) (string, error) {
	var userFileFullPath string
	// 1 - get user home directory
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	// 2 - build absolute path
	userFileFullPath = filepath.Join(userHome, relPath)

	// handle success
	return userFileFullPath, nil
}

// Description: returns a value of type T from a yaml-encoded string
//
// Example:
//
//	type FileProperty struct {
//	    Name string `json:"name"`
//	}
//
//	jsonStr := `{"name":"example.txt"}`
//	fp, err := FromJSON[FileProperty](jsonStr)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(fp.Name) // Output: example.txt
func GetVarStructFromYamlString[T any](s string) (T, error) {
	var v T
	if err := yaml.Unmarshal([]byte(s), &v); err != nil {
		return v, err
	}
	return v, nil
}

func GetVarStructFromYaml[T any](v any) (T, error) {
	var out T

	b, err := yaml.Marshal(v)
	if err != nil {
		return out, err
	}

	if err := yaml.Unmarshal(b, &out); err != nil {
		return out, err
	}

	return out, nil
}
