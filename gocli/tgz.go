package gocli

import "fmt"

func ManageTgz(filepath string) (string, error) {
	return fmt.Sprintf("manage a Go Tgz file onto the same OS:FS:%s", filepath), nil
}
