/*
Copyright © 2025 AB TRANSITION IT abtransitionit@hotmail.com

This file defines the interface any xx must implements.
*/

package errorx

// Name: Unwrapper
// Description: interface for types that can unwrap an underlying error.
// Note:
// - This is the same interface defined in the standard library's errors package.
// - With this: custom error types become compatible with errors.Is() and errors.As()
// - This means you can use the standard library's functions to check for a specific error type anywhere in your code, regardless of how many times it has been wrapped. This is a foundational principle of modern Go error handling.
type Unwrapper interface {
	Unwrap() error
}

// Name: Stacker
// Description: interface for types that contain a captured stack trace.
// Note:
// - This allows other parts of the application to check if an error has a stack trace and to access it if needed.
type Stacker interface {
	StackTrace() []uintptr
}
