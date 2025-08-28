package gocli

import "fmt"

func ManageExe(filepath string) (string, error) {
	// manage a Go Exe file onto the same OS:FS -  mean :renaming it AND move it
	return fmt.Sprintf("manage a Go Exe file onto the same OS:FS:%s", filepath), nil
}
