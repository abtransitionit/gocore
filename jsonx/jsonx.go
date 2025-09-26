package jsonx

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/abtransitionit/gocore/logx"
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

// Name: DisplayVpsDetail
//
// Description: gets a VPS:detail or a VPS:detail:field according to field.
// Returns
// - an error instead of exiting, so the caller can handle it.
func GetFilteredJson(ctx context.Context, logger logx.Logger, jsonData Json, field string) (Json, error) {
	// 1 - check parameter
	if jsonData == nil {
		return nil, fmt.Errorf("no json provided")
	}

	// 2 - Apply optional field filtering
	if field != "" {
		val, ok := GetField(jsonData, field)
		if !ok {
			return nil, fmt.Errorf("field %s not found in VPS detail", field)
		}
		// wrap in jsonx.Json to keep consistent type
		jsonData = Json{field: val}
	}

	return jsonData, nil
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
