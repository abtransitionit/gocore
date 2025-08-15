package errorx

import (
	"fmt"
	"runtime"
	"strings"
)

// Name: errorWithStack
// Description: represent an error that includes a stack trace for enhanced debugging.
// Notes:
//   - By implementing the `Unwrap()` method, it becomes compatible with Go's
//     error wrapping and inspection functions like `errors.Is()` and `errors.As()`.
type errorWithStack struct {
	err    error
	stack  []uintptr
	frames *runtime.Frames
}

// Name: Error
// Return:
// - string: the formatted error message, including the stack trace
// Note:
//   - It fulfills the `error` interface.
func (e *errorWithStack) Error() string {
	var sb strings.Builder
	sb.WriteString(e.err.Error())
	sb.WriteString("\nStack trace:\n")

	// Reset frames to iterate from the beginning.
	e.frames = runtime.CallersFrames(e.stack)
	for {
		frame, more := e.frames.Next()
		sb.WriteString(fmt.Sprintf("\t%s:%d %s\n", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return sb.String()
}

// Unwrap returns the original error, allowing for error chaining.
// This is essential for compatibility with Go 1.13+ error wrapping.
func (e *errorWithStack) Unwrap() error {
	return e.err
}

// WithStack wraps an error with a stack trace.
// It should be used when you want to add context and a stack trace to an existing error.
// The stack trace is captured at the point where this function is called.
//
// Notes:
//   - It returns nil if the original error is nil.
//   - This approach is more idiomatic as it uses a custom error type and
//     implements the `Unwrap()` method, making it compatible with `errors.Is()`
//     and `errors.As()`.
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	// We capture 10 stack frames, starting 2 frames up from this function.
	// This ensures the stack trace points to the calling code.
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)

	return &errorWithStack{
		err:   err,
		stack: pc[:n],
	}
}

// WithNoStack creates a new error without a stack trace.
// This is a convenience function that wraps `fmt.Errorf`, providing a clear and
// intentional way to create errors that don't need a stack trace, such as
// domain-level or logical errors.
//
// Notes:
//   - This is the preferred way to create a new error within this package, as it
//     provides a consistent naming scheme alongside `WithStack`.
func WithNoStack(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
