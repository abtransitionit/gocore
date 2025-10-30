package viperx

import (
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

func (c *Config) Print() error {
	confMap := make(map[string]any)
	keys := c.AllKeys()
	sort.Strings(keys)

	for _, key := range keys {
		confMap[key] = c.Get(key)
	}

	b, err := yaml.Marshal(confMap)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	fmt.Println(string(b))
	return nil
}
