package yamlx

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// Description: renders a file (Yaml, Json, etc) in memory (as []byte) using a contextual place holder structure and converts it back byte[]
//
// Usage Example:
//
//	varsPlaceholder := map[string]map[string]string{
//	    "XXX": {
//	        "aaa": repoVersion,
//	        "bbb": pkgType,
//	    },
//	}
//
// renderedFileAsBute, err := ResolveTplConfig(templatedYamlAsByte, varsPlaceholder)
func resolveTplConfig(tpl []byte, vars any) ([]byte, error) {

	// parse
	t, err := template.New("cfg").
		Option("missingkey=error"). // fail on missing key
		Parse(string(tpl))
	if err != nil {
		return nil, fmt.Errorf("parsing the templated file: %w", err)
	}

	// execute
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return nil, fmt.Errorf("resolving the templated file: %w", err)
	}

	// handle success
	return buf.Bytes(), nil
}

// Description: renders a file (Yaml, Json, etc) in memory (as []byte) using a contextual place holder structure and converts it back byte[]
func LoadTplFile(embeddedTplFile []byte, fallbackPath string, vars any) ([]byte, error) {
	source := embeddedTplFile

	// Priority: If a external path is provided, it overrides the embedded file
	if fallbackPath != "" {
		data, err := os.ReadFile(fallbackPath)
		if err != nil {
			return nil, fmt.Errorf("reading external template %q: %w", fallbackPath, err)
		}
		source = data
	}

	// Safety check: ensure we actually have content to render
	if len(source) == 0 {
		return nil, fmt.Errorf("template source is empty (no embedded data or valid fallback path)")
	}

	return resolveTplConfig(source, vars)
}

// todo
// template.FuncMap{
//     "env": os.Getenv,
// }
// kubernetesVersion: {{ env "K8S_VERSION" }}
