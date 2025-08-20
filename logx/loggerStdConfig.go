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

	const paddingWidth = 30     // fixed width for file:line field (alignment)
	const showTimestamp = false // toggle: keep or drop timestamp

	// Split log line into whitespace-separated fields
	fields := strings.Fields(line)
	if len(fields) == 0 {
		// fallback: empty line → print as-is
		return w.out.Write(p)
	}

	// --- STEP 1: handle timestamp ---
	// By convention, first 2 fields are: "YYYY/MM/DD" and "HH:MM:SS"
	var timestamp []string
	if showTimestamp && len(fields) >= 2 {
		timestamp = fields[0:2] // keep date + time
		fields = fields[2:]
	} else if !showTimestamp && len(fields) >= 2 {
		fields = fields[2:] // drop date + time
	}

	// --- STEP 2: handle file path + line number ---
	// First remaining field looks like: "path/to/file.go:42:"
	if len(fields) > 0 {
		fileField := fields[0]
		colonIndex := strings.LastIndex(fileField, ":")
		if colonIndex != -1 {
			// split path and line number
			filePath := fileField[:colonIndex]
			lineNumber := fileField[colonIndex:]

			// apply optional path formatter (e.g. shorten dirs)
			if w.pathFormatter != nil {
				filePath = w.pathFormatter(filePath)
			}

			// pad to fixed width so INFO messages align nicely
			padded := fmt.Sprintf("%-*s", paddingWidth, filePath+lineNumber)
			fields[0] = padded
		}
	}

	// --- STEP 3: rebuild log line ---
	newFields := append(timestamp, fields...)
	newLine := strings.Join(newFields, " ") + "\n"

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
