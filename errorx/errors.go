/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

gocore/errorx/errors.go implements each method of the interface.
*/

package errorx

import (
	"fmt"
	"runtime"
	"strings"
)

// The default number of stack frames to capture.
const defaultStackDepth = 10

// Name: errorWithStack
// Description: represents an error that includes a stack trace.
// Notes:
// - It implements the standard library's `Unwrap()` method and our `Stacker` interface.
type errorWithStack struct {
	msg   string
	err   error
	stack []uintptr
}

// Name: errorWithNoStack
// Description: represents an error that do not includes a stack trace.
// Notes:
// - It implements the standard library's `Unwrap()` method and our `Stacker` interface.
type errorWithNoStack struct {
	msg string
	err error
}

// Name: Error
// Return:
// - string: the formatted error message, including the stack trace.
func (e *errorWithStack) Error() string {
	var sb strings.Builder
	// If there's a wrapped error, include its message first.
	if e.err != nil {
		sb.WriteString(e.err.Error())
		sb.WriteString(": ")
	}
	sb.WriteString(e.msg)
	sb.WriteString("\n")

	// Append the stack trace.
	frames := runtime.CallersFrames(e.stack)
	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("\t%s:%d %s\n", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return sb.String()
}

// Name: Unwrap
// Return:
// - error: the original wrapped error.
// Notes:
// - This is essential for compatibility with errors.Is() and errors.As().
func (e *errorWithStack) Unwrap() error {
	return e.err
}

// Name: StackTrace
// Return:
// - []uintptr: the captured stack trace.
// Notes:
// - It fulfills the Stacker interface.
func (e *errorWithStack) StackTrace() []uintptr {
	return e.stack
}

func (e *errorWithNoStack) Error() string {
	var sb strings.Builder
	if e.err != nil {
		sb.WriteString(e.err.Error())
		sb.WriteString(": ")
	}
	sb.WriteString(e.msg)
	return sb.String()
}

func (e *errorWithNoStack) Unwrap() error {
	return e.err
}

// Name: NewWithNoStack
// Description: creates a new error with a message but no stack trace.
func NewWithNoStack(format string, a ...any) error {
	return &errorWithNoStack{
		msg: fmt.Sprintf(format, a...),
	}
}

// Name: WrapWithNoStack
// Description: wraps an existing error with a new message, without adding a stack trace.
func WrapWithNoStack(err error, format string, a ...any) any {
	if err == nil {
		return nil
	}
	return &errorWithNoStack{
		msg: fmt.Sprintf(format, a...),
		err: err,
	}
}

// Name: captureStack
// Description: captures a stack trace at the current position.
func captureStack() []uintptr {
	pc := make([]uintptr, defaultStackDepth)
	n := runtime.Callers(2, pc)
	return pc[:n]
}

// Name: New
// Description: creates a new error with a message and a captured stack trace.
// Notes:
//   - It is the primary way to generate a new root error in the application.
func New(format string, a ...any) error {
	return &errorWithStack{
		msg:   fmt.Sprintf(format, a...),
		stack: captureStack(),
	}
}

// Name: Wrap
// Description: wraps an existing error with a new message, preserving the original error and its stack trace.
// Notes:
// - If the original error is nil, it returns nil.
// - If the original error already has a stack trace, it's preserved rather than replaced.
func Wrap(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}

	// Check if the error is already of type *errorWithStack.
	// If so, we can simply wrap it without creating a new stack trace,
	// which avoids redundant and confusing stack traces.
	if _, ok := err.(*errorWithStack); ok {
		return &errorWithStack{
			msg: fmt.Sprintf(format, a...),
			err: err,
		}
	}

	// For a new error, capture a stack trace.
	return &errorWithStack{
		msg:   fmt.Sprintf(format, a...),
		err:   err,
		stack: captureStack(),
	}
}
