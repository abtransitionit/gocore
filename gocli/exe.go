package gocli

import (
	"fmt"

	"github.com/abtransitionit/gocore/logx"
)

func ManageExe(logger logx.Logger, goCli GoCli, filepath string) (string, error) {
	// logger.Debugf("🌐 Cli: %s:type:Exe:%s : now sudo copy %s to folder /usr/local/bin with name xxx", goCli.Name, filePath, filePath)
	// manage a Go Exe file onto the same OS:FS -  mean :renaming it AND move it
	return fmt.Sprintf("manage a Go Exe file onto the same OS:FS:%s", filepath), nil
}
