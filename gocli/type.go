package gocli

type GoCli struct {
	Name    string
	Version string
	Url     string
}

type MapGoCli map[string]GoCli
