// file: gocore/property/core.go
package property

import (
	"fmt"
	"strings"
)

// Name: PropertyHandler
//
// Description: retrieves a system property.
type PropertyHandler func(...string) (string, error)

// Name: GetProperty
//
// Description: retrieves a property from the os.
func GetProperty(property string, params ...string) (string, error) {

	// get function that manages that property
	fnHandler, ok := coreProperties[property]
	if !ok {
		return "", fmt.Errorf("unknown property requested: %s", property)
	}

	// play that function and get it output
	output, err := fnHandler(params...)
	if err != nil {
		return "", fmt.Errorf("error getting %s: %w", property, err)
	}

	return strings.TrimSpace(output), nil
}
