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
	"runtime"
	"strings"
)

// Name: StdLoggerConfig
// Description: holds the configuration for GO standard logger driver and custom configuration.
type StdLoggerConfig struct {
	Out           io.Writer
	Prefix        string
	Flag          int
	ShowTimestamp bool // keep or drop timestamp
	ShowFunc      bool // include or hide function name
	UseShortPath  bool // use short path via pathFormatter or full path
	PathFormatter func(string) string
}

// Name: pathWriter
//
// Description: a type that represents a writer that formats the path of the log message
type pathWriter struct {
	// out    io.Writer
	Config *StdLoggerConfig
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

	var paddingWidth = 60                      // fixed width for file:line field (alignment)
	var showTimestamp = w.Config.ShowTimestamp // toggle: keep or drop timestamp
	var showFunc = w.Config.ShowFunc           // toggle: include or hide function name
	var useShortPath = w.Config.UseShortPath   // toggle: use pathFormatter (short path) or full path

	// Split log line into whitespace-separated fields
	fields := strings.Fields(line)
	if len(fields) == 0 {
		// fallback: empty line → print as-is
		return w.Config.Out.Write(p)
	}

	// --- STEP 1: handle timestamp ---
	var timestamp []string
	if showTimestamp && len(fields) >= 2 {
		timestamp = fields[0:2] // keep date + time
		fields = fields[2:]
	} else if !showTimestamp && len(fields) >= 2 {
		fields = fields[2:] // drop date + time
	}

	// --- STEP 2: handle file path + line number ---
	if len(fields) > 0 {
		fileField := fields[0]
		colonIndex := strings.LastIndex(fileField, ":")
		if colonIndex != -1 {
			filePath := fileField[:colonIndex]
			lineNumber := fileField[colonIndex:]

			// --- optionally shorten path ---
			if useShortPath && w.Config.PathFormatter != nil {
				filePath = w.Config.PathFormatter(filePath)
			}

			// pad to fixed width so INFO messages align nicely
			padded := fmt.Sprintf("%-*s", paddingWidth, filePath+lineNumber)
			fields[0] = padded

			// --- optionally add function name ---
			if showFunc {
				if pc, _, _, ok := runtime.Caller(5); ok {
					if fn := runtime.FuncForPC(pc); fn != nil {
						funcName := fn.Name()
						fields = append([]string{fields[0], funcName}, fields[1:]...)
					}
				}
			}
		}
	}

	// --- STEP 3: rebuild log line ---
	newFields := append(timestamp, fields...)
	newLine := strings.Join(newFields, " ") + "\n"

	return w.Config.Out.Write([]byte(newLine))
}

// Name: NewStdProdConfig
// Description: creates a production configuration for GO standard logger driver.
func NewStdProdConfig() StdLoggerConfig {
	return StdLoggerConfig{
		Out:    os.Stdout,
		Prefix: "",
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	}
}

// Name: NewStdDevConfig
// Description: creates a development configuration for GO standard logger driver.
func NewStdDevConfig() StdLoggerConfig {
	// define a custom writer
	devOut := &pathWriter{
		Config: &StdLoggerConfig{
			Out:           os.Stdout,
			ShowTimestamp: false,
			ShowFunc:      false,
			UseShortPath:  true,
			PathFormatter: pathFormatter,
		},
	}
	return StdLoggerConfig{
		Out:    devOut, // use the custom writer
		Prefix: "",
		Flag:   log.LstdFlags | log.Llongfile, // date/time + caller
	}
}
