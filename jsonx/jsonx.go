package jsonx

import (
	"encoding/json"
	"fmt"

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
