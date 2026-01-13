package cilium

import (
	"fmt"

	"github.com/abtransitionit/gocore/tpl"
)

func getConfig(ciliumConf CiliumConf) (string, error) {

	// define the structure that holds the vars that will be used to resolve the templated file
	ciliumConfigFileTplVar := CiliumConf{
		K8sPodCidr:   ciliumConf.K8sPodCidr,
		K8sApiServer: ciliumConf.K8sApiServer,
	}

	// resolve the templated file
	CiliumConfigFile, err := tpl.ResolveTplConfig(configFileTpl, ciliumConfigFileTplVar)
	if err != nil {
		return "", fmt.Errorf("faild to resolve templated repo file, for the file: %s", configFileTpl)
	}

	// resturn the YamlString
	return CiliumConfigFile, nil

}
