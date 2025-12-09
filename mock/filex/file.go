package filex

import "gopkg.in/yaml.v3"

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
func GetVarStruct[T any](s string) (T, error) {
	var v T
	if err := yaml.Unmarshal([]byte(s), &v); err != nil {
		return v, err
	}
	return v, nil
}
