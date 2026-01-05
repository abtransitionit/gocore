package helm

import (
	"fmt"

	"github.com/abtransitionit/gocore/filex"
)

// -----------------------------------------
// ------ get YAML file into a var     -----
// -----------------------------------------

func getYaml(hostName string) (*MapYaml, error) {
	// 1 - get local auto cached (embedded) file into a struct
	yamlAsStruct, err := filex.LoadYamlIntoStruct[MapYaml](yamlListRepo)
	if err != nil {
		return nil, fmt.Errorf("%s > loading config: %w", hostName, err)
	}
	// handle success
	return yamlAsStruct, nil
}
