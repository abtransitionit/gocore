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