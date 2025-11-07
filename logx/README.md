# Intro
This package defines library to manage Logging. The purpose of logging is
  - to record **events** (aka. **logs**) that occur in a **software system**
  - essential for:
    - **Debugging**: Identifying the root cause of issues.
    - **Monitoring**: Tracking application health, performance, and usage.
    - **Auditing**: Providing a record of significant actions for **security** or **compliance** purposes.

**Logs** are often categorized by severity levels, such as `DEBUG`, `INFO`, `WARN`, and `ERROR`, to help developers filter and prioritize information and so action.

# The code
## `interface.go`

Defines the logging methods used by the GALI (**G**lobal **A**pplication **L**ogger **I**nstance)

## `loggerXXX.go`

Defines an implementation for a specific driver. e.g.
- `loggerZap.go` for the `Uber's zap driver`
- `loggerStd.go` for the `Go standard logging driver`


## `loggerXXXConfig.go`

configure a specific driver for all available app env. e.g.
- `loggerZapConfig.go` configure the driver to display differently in `dev` and `prod`


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
### use it directly
```go
logx.Info(...)
```
### Get it to pass it to a function
```go
func testPhase() {
	// Get the global logger instance.
	log := logx.GetLogger()

	// Run the entire sequence using the new method.
	if err := mySequence.Run(log); err != nil {
		log.ErrorWithNoStack(err, "Workflow execution failed.")
		return
	}
}

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
## Add more useful methods to interface
- add more methods to the interface to be ilmplemented. example:
```go
logx.Infof("installing %s", cliConf.Name)
logx.Errorf("Error: %v", err)
```

## Integrate package `errox` to `logx`
A logger should be **smart enough** to detect and handle **custom errors**, automatically including the stack trace when it's available.


Let's move on to the next major step: integrating your new `errorx` package with your `logx` package. This is where the true power of your production-grade design will be realized. Your logger should be smart enough to detect and handle your custom errors, automatically including the stack trace when it's available.

-----

### Step 1: Define a Logging Function for Errors

First, we need to add a new function to your `logx` interface and concrete implementations. This function, let's call it `LogErr`, will be specifically designed to handle errors. It will check if the error implements the `errorx.Stacker` interface and, if so, log the stack trace along with the error message.

#### a. Update `interface.go`

Add the `LogErr` method to your `Logger` interface.

```go
// interface.go
package logx

type Logger interface {
	Info(format string, v ...any)
	Error(format string, v ...any)
	// LogErr logs an error, including a stack trace if available.
	LogErr(err error, format string, v ...any)
}
```

This ensures that every logger implementation you create (e.g., `stdLogger`, `zapLogger`) must handle errors in this new, standardized way.

-----

### Step 2: Implement `LogErr` in Your Loggers

Now, you need to add the implementation of `LogErr` to your `stdLogger` and `zapLogger` types. This is where you'll use the `errorx` utilities we just created.

#### a. `stdlogger.go`

Modify `stdlogger.go` to use the `errorx.FormatStack` function to print the stack trace.

```go
// stdlogger.go
package logx

import (
	"log"
	"os"
	"strings"

	"github.com/abtransitionit/gocore/errorx"
)

// NewStdLogger returns a new Logger instance based on the provided configuration.
func NewStdLogger(config StdLoggerConfig) Logger {
	return &stdLogger{
		logger: log.New(config.Out, config.Prefix, config.Flag),
	}
}

func (s *stdLogger) Info(format string, v ...any) {
	s.logger.Printf("INFO: "+format, v...)
}

func (s *stdLogger) Error(format string, v ...any) {
	s.logger.Printf("ERROR: "+format, v...)
}

// LogErr logs an error message and its stack trace if it exists.
func (s *stdLogger) LogErr(err error, format string, v ...any) {
	var sb strings.Builder
	sb.WriteString("ERROR: ")
	sb.WriteString(fmt.Sprintf(format, v...))
	sb.WriteString(": ")
	sb.WriteString(err.Error())
	sb.WriteString("\n")

	if stack := errorx.GetStack(err); stack != nil {
		sb.WriteString(errorx.FormatStack(stack))
	}

	s.logger.Println(sb.String())
}
```

#### b. `zaplogger.go`

Modify `zaplogger.go` to leverage `zap`'s structured logging capabilities. Instead of a long string, you'll log the stack trace as a dedicated field.

```go
// zaplogger.go
package logx

import (
	"fmt"
	"go.uber.org/zap"
	"github.com/abtransitionit/gocore/errorx"
)

// NewZapLogger returns a new Logger instance based on a custom zap.Config.
func NewZapLogger(config zap.Config) Logger {
	l, _ := config.Build()
	return &zapLogger{
		logger: l,
	}
}

// Info logs a message with the zap logger at Info level.
func (z *zapLogger) Info(format string, v ...any) {
	z.logger.Info(fmt.Sprintf(format, v...))
}

// Error logs a message with the zap logger at Error level.
func (z *zapLogger) Error(format string, v ...any) {
	z.logger.Error(fmt.Sprintf(format, v...))
}

// LogErr logs an error message with the zap logger, including the stack trace as a structured field.
func (z *zapLogger) LogErr(err error, format string, v ...any) {
	fields := []zap.Field{
		zap.Error(err),
		zap.String("message", fmt.Sprintf(format, v...)),
	}

	if stack := errorx.GetStack(err); stack != nil {
		fields = append(fields, zap.String("stack_trace", errorx.FormatStack(stack)))
	}

	z.logger.Error(err.Error(), fields...)
}
```

-----

### Step 3: Use the New Logging Function

Now, any part of your application that needs to log an error can use the new `logx.LogErr` function. For example, your `gotest` application could now handle errors like this:

```go
// main.go (simplified example)
package main

import (
	"github.com/abtransitionit/gocore/errorx"
	"github.com/abtransitionit/gocore/logx"
)

func main() {
	logx.Init() // Initializes logger based on env vars

	err := someFunctionThatReturnsAnError()
	if err != nil {
		// The logger automatically checks for a stack trace.
		logx.LogErr(err, "An error occurred while running the task")
	}
}

func someFunctionThatReturnsAnError() error {
	// Create a new error with a stack trace.
	return errorx.New("failed to connect to the database")
}
```

This final step creates a powerful, unified logging and error management system. Your applications will now generate rich, debuggable log messages in development and clean, structured ones in production, all with a consistent and easy-to-use API.


# Exemple usage
```go
package main

import (
	"errors"
	"time"

	"github.com/abtransitionit/gocore/logx"
	"go.uber.org/zap"
)

func main() {
	// Choose dev or prod config
	config := zap.NewProductionConfig()
	logger := logx.NewZapLogger(config)

	// Simple info log
	logger.Info("Starting application...")

	// Formatted info log
	appVersion := "1.0.0"
	logger.Infof("Application version: %s", appVersion)

	// Structured info log
	logger.Infow("User logged in",
		"userID", 42,
		"role", "admin",
		"time", time.Now(),
	)

	// Simple error log
	logger.Error("Something went wrong")

	// Formatted error log
	err := errors.New("connection failed")
	logger.Errorf("Failed to connect: %v", err)

	// Structured error log
	logger.Errorw("Database error",
		"error", err,
		"retry", true,
	)

	// Error with stack trace
	logger.ErrorWithStack(err, "Critical failure in module %s", "database")

	// Error without stack trace
	logger.ErrorWithNoStack(err, "Non-critical failure in module %s", "cache")
}

```


# The process

1. `GetLogger()` is called

	```go
	if globalLogger == nil {
			Init()
	}
	```

	Since `globalLogger` is nil (first call), it calls `Init()`.


1. `Init()` runs

	```go
	once.Do(InitLogger)
	```

	So it calls `InitLogger()` exactly once.


1. `InitLogger()` calls `NewLogger()`

	```go
	globalLogger = NewLogger()
	```


1. `NewLogger()` executes

	If **no env vars are set**:

	* `appLogDriver := os.Getenv("APP_LOG_DRIVER")` → `""`
	* `appEnv := os.Getenv("APP_ENV")` → `""`

	So: First switch: `APP_LOG_DRIVER`

	```go
	switch appLogDriver {
	case "zap":
			...
	default:
	```

    Because it's empty, the **default case** is used: Use the **standard logger** (not zap).


    Inside default case: check environment `APP_ENV`

	```go
	if appEnv == "prod" {
			return NewStdLogger(NewStdProdConfig())
	}
	return NewStdLogger(NewStdDevConfig())
	```

	* `appEnv` is empty (`""`), so **not prod** → use `NewStdDevConfig()`.


Final Result (when no env vars is set) **You get:** 👉 `NewStdLogger(NewStdDevConfig())`

So the logger will be:

* **Standard Go logger**, NOT zap.
* **Development configuration**, NOT production mode.
* **Only initialized once** due to `once.Do()`.

Here is the **full decision matrix** covering every possible combination of
`APP_LOG_DRIVER` × `APP_ENV`.

This shows exactly which logger is created in each case.

---

# ✅ Full Decision Matrix

| `APP_LOG_DRIVER` | `APP_ENV`     | Selected Logger Type | Selected Config | Final Result                       |
| ---------------- | ------------- | -------------------- | --------------- | ---------------------------------- |
| *empty*          | *empty*       | Standard Logger      | Dev config      | `NewStdLogger(NewStdDevConfig())`  |
| *empty*          | `prod`        | Standard Logger      | Prod config     | `NewStdLogger(NewStdProdConfig())` |
| *empty*          | anything else | Standard Logger      | Dev config      | `NewStdLogger(NewStdDevConfig())`  |
| `zap`            | *empty*       | Zap Logger           | Dev config      | `NewZapLogger(NewDevConfig())`     |
| `zap`            | `prod`        | Zap Logger           | Prod config     | `NewZapLogger(NewProdConfig())`    |
| `zap`            | anything else | Zap Logger           | Dev config      | `NewZapLogger(NewDevConfig())`     |
| anything else    | *empty*       | Standard Logger      | Dev config      | `NewStdLogger(NewStdDevConfig())`  |
| anything else    | `prod`        | Standard Logger      | Prod config     | `NewStdLogger(NewStdProdConfig())` |
| anything else    | anything else | Standard Logger      | Dev config      | `NewStdLogger(NewStdDevConfig())`  |

---

# ✅ Summary Logic (human-readable)

* **APP_LOG_DRIVER = "zap"** → always use Zap

  * `APP_ENV = "prod"` → Zap + Prod config
  * otherwise → Zap + Dev config

* **APP_LOG_DRIVER anything else (including empty)** → use Std logger

  * `APP_ENV = "prod"` → Std + Prod config
  * otherwise → Std + Dev config

