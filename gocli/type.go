package gocli

import "fmt"

type GoCli struct {
	Name    string
	Version string
	Url     string
	OsName  string
	Doc     []string
}

type SliceGoCli []GoCli
type MapGoCli map[string]GoCli

func GetOsName(cli GoCli) (string, error) {
	// check parameters
	if cli.OsName == "" {
		return "", fmt.Errorf("no GO cli provided")
	}

	// check if cli if of type Exe
	if cli.OsName == "exe" {
		return cli.Name, nil
	}

	return cli.OsName, nil
}
