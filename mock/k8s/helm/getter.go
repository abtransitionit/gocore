package helm

import (
	"fmt"

	"github.com/abtransitionit/gocore/filex"
)

func getYaml(hostName string) (*MapYaml, error) {
	// 1 - get local auto cached (embedded) file into a struct
	yamlAsStruct, err := filex.LoadYamlIntoStruct[MapYaml](yamlList)
	if err != nil {
		return nil, fmt.Errorf("%s > loading config: %w", hostName, err)
	}
	// handle success
	return yamlAsStruct, nil
}
