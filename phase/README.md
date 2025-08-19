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
- a **tier** is a set of phase that runs concurently
  - a tier can start only when all phases of the previous tier have finished running.



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



# The execute function
This is the cornerstone of the process that execute all phases of a worflow:

- From a set of phases (ie. a worflow):

| ID  | PHASE     | DESCRIPTION                                                        | DEPENDENCIES |
|-----|-----------|--------------------------------------------------------------------|--------------|
| 1   | checklist | check VMs are SSH reachable.                                       | none         |
| 2   | cpluc     | provision LUC CLI                                                  | none         |
| 3   | dapack1   | provision standard/required/missing OS CLI (via dnfapt packages).  | [upgrade]    |
| 4   | dapack2   | provision OS dnfapt package(s) on VM(s).                           | [upgrade]    |
| 5   | gocli     | provision Go toolchain                                             | [dapack1]    |
| 6   | linger    | Allow non-root user to run OS services.                            | [dapack1]    |
| 7   | path      | configure OS PATH envvar.                                          | [dapack1]    |
| 8   | rc        | Add a line to non-root user RC file.                               | [dapack1]    |
| 9   | service   | configure OS services on Kind VMs.                                 | [dapack1]    |
| 10  | show      | display the desired KIND Cluster's configuration                   | none         |
| 11  | upgrade   | provision OS nodes with latest dnfapt packages and repositories.   | [cpluc]      |

- It creates a set of tiers:

| Tier | PHASE     | DESCRIPTION                                                        | DEPENDENCIES |
|------|-----------|--------------------------------------------------------------------|--------------|
| 1    | checklist | check VMs are SSH reachable.                                       | none         |
| 1    | **cpluc**     | provision LUC CLI                                                  | none         |
| 1    | show      | display the desired KIND Cluster's configuration                   | none         |
| 2    | upgrade   | provision OS nodes with latest dnfapt packages and repositories.   | [cpluc]      |
| 3    | dapack1   | provision standard/required/missing OS CLI (via dnfapt packages).  | [upgrade]    |
| 3    | dapack2   | provision OS dnfapt package(s) on VM(s).                           | [upgrade]    |
| 4    | gocli     | provision Go toolchain                                             | [dapack1]    |
| 4    | **linger**    | Allow non-root user to run OS services.                            | [dapack1]    |
| 4    | path      | configure OS PATH envvar.                                          | [dapack1]    |
| 4    | rc        | Add a line to non-root user RC file.                               | [dapack1]    |
| 4    | service   | configure OS services on Kind VMs.                                 | [dapack1]    |

- some phases of the worflow can be skipped thus building a **filtered** set of tiers.
- suppose we want to **skip the phases**: `cpluc` and `linger`. the **filtered** set of tiers would be:


| Tier | PHASE     | DESCRIPTION                                                        | DEPENDENCIES |
|------|-----------|--------------------------------------------------------------------|--------------|
| 1    | checklist | check VMs are SSH reachable.                                       | none         |
| 1    | show      | display the desired KIND Cluster's configuration                   | none         |
| 2    | **upgrade**   | provision OS nodes with latest dnfapt packages and repositories.   | **[cpluc]**      |
| 3    | dapack1   | provision standard/required/missing OS CLI (via dnfapt packages).  | [upgrade]    |
| 3    | dapack2   | provision OS dnfapt package(s) on VM(s).                           | [upgrade]    |
| 4    | gocli     | provision Go toolchain                                             | [dapack1]    |
| 4    | path      | configure OS PATH envvar.                                          | [dapack1]    |
| 4    | rc        | Add a line to non-root user RC file.                               | [dapack1]    |
| 4    | service   | configure OS services on Kind VMs.                                 | [dapack1]    |

- There is a potential pbs with **upgrade** that originally depend on **cpluc**
- Now each tier is executed sequentially. this mean:
  - in a tier, all phases run **concurently**
  - For a **next** tier to start, all the phases of the **previous** tier must have finished running.
	- a phase is bind to a GO function that is wrapped into a `func() error` 
	- this func() is executed 


