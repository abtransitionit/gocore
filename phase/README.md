# Intro
This package defines the concept of **phase** and **sequence** 
- a **phase** can be basically to a function
- a **sequence** is a set of functions to be executed in order

# How it works
- define a logger for your application that imlplements 
- define the sequence of function to play in sequence
```go
var mySequence = phase.NewPhaseList(
	phase.SetPhase("Setup", setupFunc, "Prepares the environment for the build."),
	phase.SetPhase("Build", buildFunc, "Compiles the source code into a binary."),
	phase.SetPhase("Test", testFunc, "Executes unit and integration tests."),
)
```
- run the sequence
```go
if err := mySequence.Run(log); err != nil {
    log.ErrorWithNoStack(err, "Workflow execution failed.")
    return
}
```