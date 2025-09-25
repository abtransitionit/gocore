package apicli

import (
	"fmt"
	"strings"
)

// Resolve the Path with given parameters
func (e Endpoint) BuildPath(params map[string]string) (string, error) {
	path := e.Path
	for k, v := range params {
		if v == "" {
			return "", fmt.Errorf("parameter %s is empty", k)
		}
		path = strings.ReplaceAll(path, "{"+k+"}", v)
	}
	return path, nil
}
