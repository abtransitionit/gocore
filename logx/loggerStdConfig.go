// File gocore/logx/loggerStdConfig.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

defines the different config concerning the GO standard logging driver for the different env we want: dev or prod.

*/

// stdloggerconfig.go
package logx

import (
	"io"
	"log"
	"os"
)

// Name: StdLoggerConfig
// Description: holds the configuration for GO standard logger driver.
type StdLoggerConfig struct {
	Out    io.Writer
	Prefix string
	Flag   int
}

// Name: NewStdDevConfig
// Description: creates a development configuration for GO standard logger driver.
func NewStdDevConfig() StdLoggerConfig {
	return StdLoggerConfig{
		Out:    os.Stdout,
		Prefix: "",
		Flag:   log.LstdFlags | log.Lshortfile, // date/time + caller
	}
}

// Name: NewStdProdConfig
// Description: creates a production configuration for GO standard logger driver.
func NewStdProdConfig() StdLoggerConfig {
	return StdLoggerConfig{
		Out:    os.Stdout,
		Prefix: "",
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
		// Flag:   log.Ldate | log.Ltime,
	}
}
