/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

This file implements each method of the interface.
*/

package errorx

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Name: Is
// Description: wraps errors.Is.
// Notes:
// - It reports whether any error in the err's chain matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Name: As
// Description: wraps errors.As.
// Notes:
//   - It finds the first error in the err's chain that matches the type of target,
//     and if it finds one, it sets target to that error value and returns true.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Name: GetStack
// Description: finds and returns the stack trace from the error chain.
// Notes:
//   - It iterates through the error chain and returns the first stack trace it finds
//     that implements the Stacker interface. If no such error is found, it returns nil.
func GetStack(err error) []uintptr {
	if err == nil {
		return nil
	}

	var stacker Stacker
	if As(err, &stacker) {
		return stacker.StackTrace()
	}

	// If As fails, we can also check for a concrete type match for robustness.
	// This part is a safety net in case As doesn't work as expected with the interface.
	var ews *errorWithStack
	if As(err, &ews) {
		return ews.StackTrace()
	}

	return nil
}

// Name: FormatStack
// Description: formats a captured stack trace into a human-readable string.
func FormatStack(stack []uintptr) string {
	if len(stack) == 0 {
		return "No stack trace available."
	}

	var sb strings.Builder
	sb.WriteString("Stack trace:\n")

	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		// Skip frames from the errorx package itself to provide a cleaner stack trace
		// that starts at the user's code.
		if strings.HasPrefix(frame.Function, "github.com/abtransitionit/gocore/errorx") {
			if !more {
				break
			}
			continue
		}
		sb.WriteString(fmt.Sprintf("\t%s:%d %s\n", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return sb.String()
}
