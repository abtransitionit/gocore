// File gocore/logx/loggerStdConfig.go
/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

defines the different config concerning the GO standard logging driver for the different env we want: dev or prod.

*/

// stdloggerconfig.go
package logx

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Name: StdLoggerConfig
// Description: holds the configuration for GO standard logger driver.
type StdLoggerConfig struct {
	Out    io.Writer
	Prefix string
	Flag   int
	// PathFormatter func(string) string // optional, for customizing file paths

}

// Name: pathWriter
//
// Description: a type that represents a writer that formats the path of the log message
type pathWriter struct {
	out           io.Writer
	pathFormatter func(string) string
}

// Name: Write
//
// Description: that takes a full path and returns a shorter version
//
// Parameters:
// - p: []byte: the full path to be formatted
//
// Returns:
// - n: int: the number of bytes written
// - err: error: any error that occurred during the formatting
// Notes:
//   - Every time log.Logger writes something, it will call it
func (w *pathWriter) Write(p []byte) (n int, err error) {
	line := string(p)
	const paddingWidth = 30

	// Try to detect the file path (for log.Llongfile / log.Lshortfile)
	// Standard log format is:
	//   "2009/01/23 01:23:23 /full/path/to/file.go:123: message"
	// -> first 2 fields = date + time
	// -> 3rd field      = file path + line number
	fields := strings.Fields(line)
	if len(fields) < 3 {
		// fallback: write as-is if unexpected format
		return w.out.Write(p)
	}

	// Extract file path + line number (3rd field)
	fileField := fields[2]
	colonIndex := strings.LastIndex(fileField, ":")
	if colonIndex != -1 {
		filePath := fileField[:colonIndex]
		lineNumber := fileField[colonIndex:] // keep ":123"

		// PATH: apply custom formatter (e.g. shorten to last 3 dirs)
		if w.pathFormatter != nil {
			filePath = w.pathFormatter(filePath)
		}

		// PADDING: ensure fixed width alignment for all paths
		padded := fmt.Sprintf("%-*s", paddingWidth, filePath+lineNumber)
		fields[2] = padded
	}

	// Rebuild log line with aligned path + INFO
	newLine := strings.Join(fields, " ") + "\n"
	return w.out.Write([]byte(newLine))
}

// Name: NewStdProdConfig
// Description: creates a production configuration for GO standard logger driver.
func NewStdProdConfig() StdLoggerConfig {
	return StdLoggerConfig{
		Out:    os.Stdout,
		Prefix: "",
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
		// PathFormatter: pathFormatter,
	}
}

// Name: NewStdDevConfig
// Description: creates a development configuration for GO standard logger driver.
func NewStdDevConfig() StdLoggerConfig {
	devOut := &pathWriter{
		out:           os.Stdout,
		pathFormatter: pathFormatter,
	}
	return StdLoggerConfig{
		Out:    devOut,
		Prefix: "",
		Flag:   log.LstdFlags | log.Llongfile, // date/time + caller
		// PathFormatter: pathFormatter,
	}
}
