/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com
*/

package errorx

import (
	"fmt"
	"runtime"
	"strings"
)

// WithStack returns a new error that wraps the original error and adds a stack trace.
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var sb strings.Builder
	sb.WriteString(err.Error())
	sb.WriteString("\nStack trace:\n")

	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("\t%s:%d %s\n", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}

	return fmt.Errorf("%s", sb.String())
}
