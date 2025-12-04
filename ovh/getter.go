package ovh

import (
	_ "embed"
	"fmt"

	"github.com/abtransitionit/gocore/yamlx"
)

// -----------------------------------------
// ------ define file location -------------
// -----------------------------------------

//go:embed db.list.image.yaml
var yamlVpsImage []byte // cache the raw yaml file in this var

// -----------------------------------------
// ------ get YAML file --------------------
// -----------------------------------------

// ####### of vps image name and id #######

func getVpsImageList() (*VpsImgYamlList, error) {
	theYaml, err := yamlx.LoadTplYamlFileEmbed[VpsImgYamlList](yamlVpsImage, "")
	if err != nil {
		return nil, fmt.Errorf("getting YAML file in package > %w", err)
	}
	return theYaml, nil
}
