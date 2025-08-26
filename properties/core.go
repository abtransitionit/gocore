// file: gocore/properties/core.go
package properties

import (
	"fmt"
	"strings"
)

// Name: GetPropertyLocal
//
// Description: retrieves a property from the core set or Linux-specific set if running Linux.
// GetPropertyLocal retrieves a property from gocore or golinux (if Linux).

func GetPropertyLocal(property string, params ...string) (string, error) {
	fnPropertyHandler, ok := coreProperties[property]
	if !ok {
		return "", fmt.Errorf("❌ unknown property requested: %s", property)
	}

	output, err := fnPropertyHandler(params...)
	if err != nil {
		return "", fmt.Errorf("❌ error getting %s: %w", property, err)
	}
	return strings.TrimSpace(output), nil
}
