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

	const lenValue = 90

	type Entry struct {
		Path  string
		Type  string
		Value string
	}

	var entries []Entry

	detectType := func(value interface{}) string {
		switch v := value.(type) {
		case []interface{}:
			if len(v) == 0 {
				return "[]any"
			}
			allStrings := true
			for _, elem := range v {
				if _, ok := elem.(string); !ok {
					allStrings = false
					break
				}
			}
			if allStrings {
				return "[]string"
			}
			return "[]any"
		default:
			return "string"
		}
	}

	formatValue := func(value interface{}) string {
		if value == nil {
			return ""
		}
		var s string
		switch v := value.(type) {
		case string:
			s = v
		case []interface{}:
			parts := make([]string, len(v))
			for i, elem := range v {
				parts[i] = fmt.Sprintf("%v", elem)
			}
			s = strings.Join(parts, "; ")
		default:
			b, err := yaml.Marshal(v)
			if err != nil {
				s = fmt.Sprintf("%v", v)
			} else {
				s = string(b)
				s = strings.ReplaceAll(s, "\n", "; ")
				s = strings.ReplaceAll(s, "  ", "")
			}
		}
		if len(s) > lenValue {
			s = s[:lenValue-3] + "..."
		}
		return s
	}

	var walk func(prefix string, value interface{})
	walk = func(prefix string, value interface{}) {
		switch v := value.(type) {
		case map[string]interface{}:
			for k, val := range v {
				newPrefix := k
				if prefix != "" {
					newPrefix = prefix + "." + k
				}
				walk(newPrefix, val)
			}
		default:
			entries = append(entries, Entry{
				Path:  prefix,
				Type:  detectType(v),
				Value: formatValue(v),
			})
		}
	}

	walk("", c.AllSettings())

	// simple insertion sort by Type
	for i := 1; i < len(entries); i++ {
		j := i
		for j > 0 && entries[j-1].Type > entries[j].Type {
			entries[j-1], entries[j] = entries[j], entries[j-1]
			j--
		}
	}

	// build tab-separated string
	var sb strings.Builder
	sb.WriteString("Var Name\tType\tValue\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%s\t%s\t%s\n", e.Path, e.Type, e.Value))
	}

	return sb.String(), nil
}
