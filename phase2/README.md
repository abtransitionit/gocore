# Purpose
- A framework to execute complex actions on different set of `VM` (**V**irtual **M**acchines)
- Actions and confifugration are define in a `YAML` txt file

# Terminology
## Workflow
### Formal definition
```go
type Workflow struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Phases      map[string]Phase `yaml:"phases"`
}
```
### definition
a named set of `phases`

## Phase
### Formal definition
```go
type Phase struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Fn          string   `yaml:"fn"`
	Dependency  []string `yaml:"dependency,omitempty"`
	Param       []string `yaml:"param,omitempty"`
	Node        string   `yaml:"node,omitempty"`
}
```

### Definition 
A function that is:
  - executed on a named set of **nodes**,
  - passed **parameters**
  - triggered once all **dependent** phases have completed


## Tier
the field `dependency` of a `phase` **implies** the concept of **Tier**

### Definition by example
suppose the following workflow of 4 phases:
```yaml
name: tiern
description: create a KBE (Kubernetes Easy) cluster
phases:
  APhase:
    description: run alone
    node: all
    fn: vm.CheckVmSshAccess
    Dependency: []

  BPhase:
    description: run concurently with C
    node: all
    fn: luc.DeployLuc
    Dependency: []

  CPhase:
    description: run concurently with B
    node: all
    fn: luc.DeployLuc
    Dependency: []

  DPhase:
    description: run after B and C finish
    node: all
    fn: luc.DeployLuc
    Dependency: []

```

that can be sumarized in the following table of 4 **independent** phases:
```
┌────┬────────┬──────┬─────────────────────┬───────┐
│ ID │ PHASE  │ NODE │ FN                  │ PARAM │
├────┼────────┼──────┼─────────────────────┼───────┤
│  1 │ APhase │ all  │ vm.CheckVmSshAccess │ none  │
│  2 │ BPhase │ all  │ luc.DeployLuc       │ none  │
│  3 │ CPhase │ all  │ luc.DeployLuc       │ none  │
│  4 │ DPhase │ all  │ luc.DeployLuc       │ none  │
└────┴────────┴──────┴─────────────────────┴───────┘
```
that also can be sumarized when we introduced the concept of **tier** as:
```
┌────┬──────┬─────┬────────┬──────┬──────────────────────────┬──────────────┐
│ ID │ TIER │ IDP │ PHASE  │ NODE │ DESCRIPTION              │ DEPENDENCIES │
├────┼──────┼─────┼────────┼──────┼──────────────────────────┼──────────────┤
│  1 │ 1    │ 1   │ APhase │ all  │ run alone                │ none         │
│  2 │ 1    │ 2   │ BPhase │ all  │ run concurently with C   │ none         │
│  3 │ 1    │ 3   │ CPhase │ all  │ run concurently with B   │ none         │
│  4 │ 1    │ 4   │ DPhase │ all  │ run after B and C finish │ none         │
└────┴──────┴─────┴────────┴──────┴──────────────────────────┴──────────────┘
```
- Phases of the same tier run in **parallel** (ie. **concurently**). 
- Next **tier** starts when the previous one complete.
- One Phase run **concurently** on all **node**
- in this case, the 4 independent phases will run **concurently**.

suppose now the same phases with **dependencies** between them:
```go
name: tiern
description: create a KBE (Kubernetes Easy) cluster
phases:
  APhase:
    description: run alone
    node: all
    fn: vm.CheckVmSshAccess
    dependency: []

  BPhase:
    description: run concurently with C
    node: all
    fn: luc.DeployLuc
    dependency:
      - APhase

  CPhase:
    description: run concurently with B
    node: all
    fn: luc.DeployLuc
    dependency:
      - APhase

  DPhase:
    description: run after B and C finish
    node: all
    fn: luc.DeployLuc
    dependency:
      - BPhase
      - CPhase

```
  

that can be sumarized in the following table of 4 **dependent** phases (in which tiers does not appear):
```
┌────┬────────┬──────┬─────────────────────┬───────┐
│ ID │ PHASE  │ NODE │ FN                  │ PARAM │
├────┼────────┼──────┼─────────────────────┼───────┤
│  1 │ APhase │ all  │ vm.CheckVmSshAccess │ none  │
│  2 │ BPhase │ all  │ luc.DeployLuc       │ none  │
│  3 │ CPhase │ all  │ luc.DeployLuc       │ none  │
│  4 │ DPhase │ all  │ luc.DeployLuc       │ none  │
└────┴────────┴──────┴─────────────────────┴───────┘
```
that also can be sumarized when we introduced the concept of **tier** as:
```
┌────┬──────┬─────┬────────┬──────┬──────────────────────────┬────────────────┐
│ ID │ TIER │ IDP │ PHASE  │ NODE │ DESCRIPTION              │ DEPENDENCIES   │
├────┼──────┼─────┼────────┼──────┼──────────────────────────┼────────────────┤
│  1 │ 1    │ 1   │ APhase │ all  │ run alone                │ none           │
│  2 │ 2    │ 1   │ BPhase │ all  │ run concurently with C   │ APhase         │
│  3 │ 2    │ 2   │ CPhase │ all  │ run concurently with B   │ APhase         │
│  4 │ 3    │ 1   │ DPhase │ all  │ run after B and C finish │ BPhase, CPhase │
└────┴──────┴─────┴────────┴──────┴──────────────────────────┴────────────────┘
```
- Phases of the same tier run in **parallel** (ie. **concurently**). 
- Next **tier** starts when the previous one complete.
- One Phase run **concurently** on all **node**
- in this case
  - there are 3 **tiers**.
  - tier 1 consist of the only phase `APhase` that will be executed first (it will be executed concurently on all node).
  - tier 2 consist of the 2 phases `BPhase` and `CPhase` that will be executed concurently after `APhase` complete (on all node).
  - tier 3 consist of the only phase `DPhase` that will be executed after all phase of **tier 2**` complete.


## function registry
### formal definition
```go
type FunctionRegistry struct {
	funcs map[string]any
}
```
### definition
a dictionnary that get an **object** from a **string**

# How it works
**step1**: define the **configuration** of a worklow (e.g. named **simple**) in a `YAML` txt file 

```yaml
workflow:
  simple:

    node:
      all:
        - o1u
        - o2a

```

**step2**: define the **workflow** named **simple** in a `YAML` txt file 

```yaml
name: simple
description: a small workflow
phases:
  NodeSsh:
    description: Check if VMs are SSH reachable
    node: all
    fn: vm.CheckVmSshAccess
    Dependency: []

  LucAgent:
    description: Copy LUC CLI agent to all VMs
    node: all
    fn: luc.DeployLuc
    Dependency:
      - NodeSsh

  OsUpgrade:
    description: Provision OS nodes with latest dnfapt packages and repositories
    node: all
    fn: dnfapt.UpgradeVmOs
    Dependency:
      - LucAgent

  AppUpdate:
    description: Provision required/missing standard dnfapt packages
    node: all
    fn: dnfapt.UpdateVmOsApp
    params:
      daPkgKey: da.pkg.required
    Dependency:
      - OsUpgrade

```

the workflow can be displayed via different view:
```
$ goluc wkf simple print --phase
┌────┬───────────┬──────┬──────────────────────┬───────┐
│ ID │ PHASE     │ NODE │ FN                   │ PARAM │
├────┼───────────┼──────┼──────────────────────┼───────┤
│  1 │ AppUpdate │ all  │ dnfapt.UpdateVmOsApp │ none  │
│  2 │ LucAgent  │ all  │ luc.DeployLuc        │ none  │
│  3 │ NodeSsh   │ all  │ vm.CheckVmSshAccess  │ none  │
│  4 │ OsUpgrade │ all  │ dnfapt.UpgradeVmOs   │ none  │
└────┴───────────┴──────┴──────────────────────┴───────┘

$ goluc wkf simple print --tier

┌────┬──────┬─────┬───────────┬──────┬─────────────────────────────────────────────────────────────────┬──────────────┐
│ ID │ TIER │ IDP │ PHASE     │ NODE │ DESCRIPTION                                                     │ DEPENDENCIES │
├────┼──────┼─────┼───────────┼──────┼─────────────────────────────────────────────────────────────────┼──────────────┤
│  1 │ 1    │ 1   │ AppUpdate │ all  │ Provision required/missing standard dnfapt packages             │ none         │
│  2 │ 1    │ 2   │ LucAgent  │ all  │ Copy LUC CLI agent to all VMs                                   │ none         │
│  3 │ 1    │ 3   │ NodeSsh   │ all  │ Check if VMs are SSH reachable                                  │ none         │
│  4 │ 1    │ 4   │ OsUpgrade │ all  │ Provision OS nodes with latest dnfapt packages and repositories │ none         │
└────┴──────┴─────┴───────────┴──────┴─────────────────────────────────────────────────────────────────┴──────────────┘
```
**step3**: execute the workflow: 
```
$goluc wkf simple run
```

At **runtime** (when executing this command)
  - The **workflow** YAML is loaded into a `Go` struct.
  - The **config** YAML is loaded into a `Go` struct.
  - Each `fn` function of a phase, according to the rule of **tier**, is:
    - resolved from the **config** then **fetch** in the library code (package and function name)
    - passed **parameters** if any (resolved from the **config**)
    - executed on the **node**

# Errors  
- At **runtime**, if any phase **fails**, it logs errors and **stops** the workflow.  

# Resolving Function
This consist of mapping a function defined as a string to a real Go function inside a package. 
- the function is added to the registry 

# Code Detail
## Tier
- a **tier** is a slice of phases (i.e. `[][]phase`)
- for each tier we globally
  - run all phase concurently (i.e. each phase in its own goroutine)
  - wait for all phases ti finish
  - collect errors
  - when any phase failed   → the code stop and return an error.
  - when all phases succeed → move to next tier.

# Todo
- skip and retain