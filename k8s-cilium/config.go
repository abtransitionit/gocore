package cilium

import (
	"bytes"
	"fmt"
	"html/template"
)

func getConfig(ciliumConf CiliumConf) (string, error) {

	// define the structure that holds the vars that will be used to resolve the templated file
	ciliumConfigFileTplVar := CiliumConf{
		K8sPodCidr:   ciliumConf.K8sPodCidr,
		K8sApiServer: ciliumConf.K8sApiServer,
	}

	// resolve the templated file
	CiliumConfigFile, err := resolveTplConfig(configFileTpl, ciliumConfigFileTplVar)
	if err != nil {
		return "", fmt.Errorf("faild to resolve templated repo file, for the file: %s", configFileTpl)
	}

	// resturn the YamlString
	return CiliumConfigFile, nil

}

func resolveTplConfig(tplFile string, vars CiliumConf) (string, error) {
	tpl, err := template.New("repo").Parse(tplFile)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, vars); err != nil {
		return "", err
	}

	return buf.String(), nil
}
