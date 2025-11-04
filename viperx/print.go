package viperx

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Description: returns the viperx instance as a string to be print
//
// Notes:
//
// - Allow to view the content of the viperx instance
func (c *Viperx) GetContentAsString() (string, error) {

	configMap := make(map[string]any)
	keys := c.AllKeys()
	sort.Strings(keys)

	for _, key := range keys {
		configMap[key] = c.Get(key)
	}

	b, err := yaml.Marshal(configMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	return string(b), nil
}

// Description: returns the viperx instance as a string to be as a table
//
// Notes:
//
// - Allow to view the content of the viperx instance
func (c *Viperx) GetContentAsTable() (string, error) {
	if c == nil || c.Viper == nil {
		return "", fmt.Errorf("Viperx instance is nil")
	}

	type Entry struct {
		Path string
		Type string
	}

	var entries []Entry

	// file
	if v := c.Get("file.customRcFileName"); v != nil {
		entries = append(entries, Entry{"file.customRcFileName", "string"})
	}
	if v := c.Get("file.binFolderPath"); v != nil {
		entries = append(entries, Entry{"file.binFolderPath", "string"})
	}

	// node
	if v := c.GetStringSlice("node.all"); len(v) > 0 {
		entries = append(entries, Entry{"node.all", "[]string"})
	}

	// da.pkg.required
	if v := c.GetStringSlice("da.pkg.required"); len(v) > 0 {
		entries = append(entries, Entry{"da.pkg.required", "[]string"})
	}

	// da.repo.node
	if v := c.Get("da.repo.node"); v != nil {
		entries = append(entries, Entry{"da.repo.node", "[]any"})
	}

	// goCli
	if v := c.Get("goCli"); v != nil {
		entries = append(entries, Entry{"goCli", "[]any"})
	}

	// service
	if v := c.Get("service"); v != nil {
		entries = append(entries, Entry{"service", "[]any"})
	}

	// envar
	if v := c.Get("envar"); v != nil {
		entries = append(entries, Entry{"envar", "[]any"})
	}

	// helm.release
	if v := c.Get("helm.release"); v != nil {
		entries = append(entries, Entry{"helm.release", "[]any"})
	}

	// cluster
	if v := c.Get("cluster"); v != nil {
		entries = append(entries, Entry{"cluster", "[]any"})
	}

	// build raw string with tab-separated columns
	var sb strings.Builder
	sb.WriteString("Name\tType\n") // header
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%s\t%s\n", e.Path, e.Type))
	}

	return sb.String(), nil
}
