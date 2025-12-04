package ovh

import (
	_ "embed"
	"fmt"

	"github.com/abtransitionit/gocore/yamlx"
)

// -----------------------------------------
// ------ define file location -------------
// -----------------------------------------

//go:embed db.list.distro.yaml
var yamlVpsDistro []byte // cache the raw yaml file in this var

//go:embed db.list.vps.yaml
var yamlVps []byte // cache the raw yaml file in this var

// -----------------------------------------
// ------ get cached YAML file -------------
// -----------------------------------------

// ####### of ovh vps distro manage by the organization #######

func getVpsDistroList() (*DistroYaml, error) {
	theYaml, err := yamlx.LoadTplYamlFileEmbed[DistroYaml](yamlVpsDistro, "")
	if err != nil {
		return nil, fmt.Errorf("getting YAML file in package > %w", err)
	}
	return theYaml, nil
}

// ####### of ovh vps manage by the organization #######

func getVpsList() (*VpsYaml, error) {
	theYaml, err := yamlx.LoadTplYamlFileEmbed[VpsYaml](yamlVps, "")
	if err != nil {
		return nil, fmt.Errorf("getting YAML file in package > %w", err)
	}
	return theYaml, nil
}
