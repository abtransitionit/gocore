package gocli

import "fmt"

func ManageTgz(filepath string) (string, error) {
	// manage a Go Tgz file onto the same OS:FS -  mean : more job: extract and then decide
	return fmt.Sprintf("manage a Go Tgz file onto the same OS:FS:%s", filepath), nil
}
