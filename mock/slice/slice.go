package slice

import "strings"

func GetStringWithSepFromSlice(ListString []string, separator string) string {
	return strings.Join(ListString, separator)
}
