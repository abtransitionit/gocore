package helm

import (
	_ "embed"
)

// -----------------------------------------
// ------ define file location -------------
// -----------------------------------------

//go:embed db.repo.yaml
var yamlListRepo []byte // automatically cache the raw yaml file in this var
