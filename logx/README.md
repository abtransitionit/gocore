# Intro
**`interface.go`**

Defines the logging methods used by the GALI (**G**lobal **A**pplication **L**ogger **I**nstance)

**`loggerXXX.go`**

Defines an implementation for a specific driver. e.g.
- `loggerZap.go` for the `Uber's zap driver`
- `loggerStd.go` for the `Go standard logging driver`


**`loggerXXXConfig.go`**

configure a specific driver for all available app env. e.g.
- `loggerZapConfig.go` configure the driver to display differently in `dev` and `prod`


# How it works
In Development mode
```
go run main.go
```
In Production mode
```
APP_LOG_DRIVER=zap APP_ENV=prod go run main.go
APP_LOG_DRIVER=zap APP_ENV=dev go run main.go
```

# zap dev Vs. zap prod
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
2025/08/14 15:52:12 logx.go:25: INFO: Test some code.
```
- prod: 
```json
{"level": "info", "time": "2025-08-14T15:52:12Z", "caller": "logx.go:25", "msg": "Test some code", "task": "play", "component": "test"}
```