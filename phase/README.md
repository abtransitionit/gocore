# Intro

This package defines the following concepts
- a **Task** is the conceptual idea of a single, atomic unit of work to be performed.
- a **Phase** 
  - is the concrete implementation of a **Task**
  - can simply be defined by a GO function (but not only)
  - is technically the following struct

	```go
	type Phase struct {
	Name        string
	Description string
	fn          PhaseFunc
	}
	```

- a **Workflow** is the the conceptual idea of a collection of one or more **tasks**
- a **Phase list** 
  - is the concrete implementation of a **Workflow**
  - can simply be defined by a list of **phases**
	- is technically the following struct

	```go
	type PhaseList []Phase
	```
- a **Step**is a synonym for a **Phase** or a **Task**.

-  a **DAG (Directed Acyclic Graph)** is the data structure we will use to represent the dependencies within a **Workflow**. It's a way to organize our **Phases**. it will allow to:
	- manage the order in which phases are run
	- manage dependencies if any between phases
	- manage any concurency run of phases

- an **adapater** that allows each phases to potentially run concurently using the package `syncx` 
- a **context** that allows to interact with long running process



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


# Implemennting a Graceful Shutdown
**without context**
 - if your application is interrupted (e.g., via Ctrl+C)
 - the workflow will stop abruptly: potentially creating data corruption or orphaned processes


**improvemnt** 
 - add a mechanism to capture the interrupt signal
 - use the context.Context to stop all running Goroutines cleanly.
 - Prevent data corruption or orphaned processes for any long-running application.


**with context**
- We add a context that listens for an interrupt signal (like Ctrl+C). 
- When the signal is received, the context will be canceled.
- Propagate the context: The cancellable context will be passed to the workflow
- the Execute method will then pass the context to the syncx.RunConcurrently function.
- the syncx.RunConcurrently function will watch the context for a cancellation signal. 
- If the context is canceled, it will stop launching new goroutines and handle any currently running ones.



