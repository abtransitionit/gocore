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

	// Try to detect the file path at the end (for log.Llongfile / log.Lshortfile)
	// Log format is: "2009/01/23 01:23:23 /full/path/to/file.go:123: message"
	// We split on spaces, last-but-one is file path with line number
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return w.out.Write(p) // fallback
	}

	// The file path is usually the 3rd field (after date + time)
	fileField := fields[2]
	// remove trailing colon and line number
	colonIndex := strings.LastIndex(fileField, ":")
	if colonIndex != -1 {
		filePath := fileField[:colonIndex]
		lineNumber := fileField[colonIndex:] // keep ":123"
		if w.pathFormatter != nil {
			filePath = w.pathFormatter(filePath)
		}
		fields[2] = filePath + lineNumber
	}

	newLine := strings.Join(fields, " ") + "\n"
	return w.out.Write([]byte(newLine))
}

// func (w *pathWriter) Write(p []byte) (n int, err error) {
// 	line := string(p)

// 	// Here we’ll replace the full path with a shorter version
// 	if w.pathFormatter != nil {
// 		line = w.pathFormatter(line)
// 	}

// 	return w.out.Write([]byte(line))
// }

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

// func pathFormatter(fullPath string) string {
// 	parts := strings.Split(filepath.ToSlash(fullPath), "/")
// 	if len(parts) <= 3 {
// 		return fullPath
// 	}
// 	return strings.Join(parts[len(parts)-3:], "/")
// }

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
