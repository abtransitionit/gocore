# Intro

This package defines library to manage **errors**. An error is an **event** that:
 - indicates something has gone wrong during the **execution of some code**
 - must be **logged**
 - can be due to 
   - a bug
   - an unexpected condition, context, ..
   - a failure of an external resource, ...
 - often carry additional information, such as a **stack trace**, to pinpoint the exact location in the code where the problem **originated**.  

The Error management is the practice of:
- **Detecting** when an error occurs.
- **Propagating** the error up the call stack to a point where it can be handled.
- **Handling** the error gracefully, which might involve 
  - **healing** before **retrying** an operation
  - **informing** all or some  users
  - **logging** the failure.

# The code

## `interface.go`

`Unwrapper`: 
- By implementing this method, custom error types become compatible with errors.Is() and errors.As(). 
- This means you can use the standard library's functions to check for a specific error type anywhere in your code, regardless of how many times it has been wrapped. This is a foundational principle of modern Go error handling.

`Stacker`: 
- Allows to create a specific contract for errors that contain a stack trace. 
- This enables to write functions that can check if an error has a stack trace and then process it without needing to know the **concrete** type of the error. 
- For example, a logging package could check for this interface and log the stack trace only if it's available.

## `erors.go`

`errorWithStack`: 
- This struct is now the single, rich error type. It encapsulates both the error message and the stack trace.

`Error()` method: 
- handles wrapped errors more gracefully, producing a more readable output like "original error: new message".

`New()` function: 
- the canonical way to create a brand new error. It always captures a stack trace, making every root error easily debuggable.

`Wrap()` function: 
- If the error you're wrapping already has a stack trace (i.e., it's an errorWithStack type), it just wraps the existing error without adding a new, redundant stack trace. This prevents your log files from having a cascade of stack traces for a single error. 
- If the error doesn't have a stack trace, Wrap will add one. This makes error wrapping highly efficient and meaningful.


## `utils.go`

- This file contains helper functions that make it easier for your organization's developers to work with your custom errors, i.e. inspecting, checking,and extracting information from errors in a consistent way.

# How it works
## `require` the module in the `go.mod` of your project
```go
require github.com/abtransitionit/gocore v1.0.0
```

## `import` the package in the project's `main.go`
```go
import "github.com/abtransitionit/gocore/logx
```

## Use the logger in the code
```go
logx.Info(...)
```
## See it in action
  
  **In Development mode**
  ```go
  // use GO std logger driver, no matter the env
  go run main.go ...
  
  // use Zap logger driver
  APP_LOG_DRIVER=zap go run main.go ...
  ```
  
  **In Production mode**
  ```go
  // use GO std logger driver, no matter the env
  go run main.go ...

// use Zap logger driver with a prod configuration 
  APP_LOG_DRIVER=zap APP_ENV=prod go run main.go ...

// use Zap logger driver with a dev configuration 
  APP_LOG_DRIVER=zap APP_ENV=dev go run main.go ..
  ```

# configuring the Zap logger driver
The code provide a simple configuration for 2 environments named `dev` and `prod` 

## zap dev Vs. zap prod
|Feature|zap.NewDevelopment()|zap.NewProduction()|
|-|-|-|
|**Log Format**|	Human-readable console output.|Structured, machine-readable JSON output.
|**Performance**|	Slower due to extra formatting for human-readability.|Optimized for speed with zero-allocation logging.
|**Log Level**|	By default, logs at the DEBUG level and above.|By default, logs at the INFO level and above.
|**Caller Info**|	Includes file and line number by default.|Includes file and line number by default.
|**Stack Traces**|	Prints stack traces on DPanic and Panic levels.|Prints stack traces on ERROR and FATAL levels.

## Example
**Log Format**:
- dev: 
```
INFO	logx/loggerZap.go:40	Test some code.
```
- prod: 
```json
{"level":"info","ts":1755256591.774873,"caller":"logx/loggerZap.go:40","msg":"Test some code."}
```

# Todo
- add more methods to the interface to be ilmplemented. example:
```go
logx.Infof("installing %s", cliConf.Name)
logx.Errorf("Error: %v", err)
```