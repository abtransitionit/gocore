package jsonx

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/pretty"
)

func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error pretty printing:", err)
		return
	}
	fmt.Println(string(b))
}

func GetField(m map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	var current interface{} = m
	for _, p := range parts {
		if m2, ok := current.(map[string]interface{}); ok {
			if val, exists := m2[p]; exists {
				current = val
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}
	return current, true
}

func PrettyPrintColor(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Colorize and indent
	prettyJSON := pretty.Color(pretty.Pretty(b), pretty.TerminalStyle)
	fmt.Println(string(prettyJSON))
}
