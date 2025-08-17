# Gocore

A foundational Go library (i.e. no `main()`) containing low level universal cross-platform reusable building blocks, and utilities for all `abtransitionit` GO projects.

---

[![Dev CI](https://github.com/abtransitionit/gocore/actions/workflows/ci-dev.yaml/badge.svg?branch=dev)](https://github.com/abtransitionit/gocore/actions/workflows/ci-dev.yaml)
[![Main CI](https://github.com/abtransitionit/gocore/actions/workflows/ci-main.yaml/badge.svg?branch=main)](https://github.com/abtransitionit/gocore/actions/workflows/ci-main.yaml)
[![LICENSE](https://img.shields.io/badge/license-Apache_2.0-blue.svg)](https://choosealicense.com/licenses/apache-2.0/)

----


# Features  
This project template includes the following components:  


|Component|Description|
|-|-|
|Licensing|Predefined open-source license (Apache 2.0) for legal compliance.|
|Code of Conduct| Ensures a welcoming and inclusive environment for all contributors.|  
|README|Structured documentation template for clear project onboarding.|  

This repository contains core packages designed to be used across all of our services, such as:

- `errorx`: A package for structured error handling.
- `logx`: A package for consistent logging.
- `filex`: Utilities for common file system operations.


---




# Getting started
To use this module (library) in your project, run:


- `go` `get` [github.com/abtransitionit/gocore](https://github.com/abtransitionit/gocore)

# Contributing  

We welcome contributions! Before participating, please review:  
- **[Code of Conduct](.github/CODE_OF_CONDUCT.md)** – Our community guidelines.  
- **[Contributing Guide](.github/CONTRIBUTING.md)** – How to submit issues, PRs, and more.  


# Contributing as developer



## updating the `go/mod`


During local development, we use the `replace` directive in `go.mod` to simplify dependency management when working with multiple interdependent Go projects.

For promotion to production:

* The `go.mod` file is committed **as-is** into all branches: feature and dev branches.
* When code reaches the `main` (or any production) branch, the CI pipeline automatically removes all `replace` directives before building and/or generating the release artifacts.

This approach ensures that:

* Developers benefit from faster iteration and easier local linking during development.
* Production releases builds always rely on published module versions, ensuring stability and reproducibility.
 
## Updating an `interface`
1. **Modify the Interface Definition**: 
    - define and/or update the method signature. 
    - This change will immediately break the build for all code that uses a type that implements this interface.
1. **Identify Implementing Types**: 
    - launch a `go vet` or `go build`
    - failures allows to identify code to update
1. **Update the Implementing Types**: 
    - For each of the types you identified
    - add/update the new method to match the updated interface signature
    - To get your code to compile quickly: providing a method stubs
        - that have the correct signature but contain minimal logic.
        - this allows you to restore a working build and then implement the full functionality later.

Here’s a polished, professional, and production-ready version of your text with consistent tone, grammar, and clarity:


## Testing Code

As the Go community, we follow a **white-box testing** approach, which allows us to test not only public but also private functions. This is achieved by:

* **Keeping test files alongside the source code**, ensuring maintainability and readability.
* **Promoting a strong unit testing culture**, where tests are considered a first-class part of the codebase.


### Go Test Files

A Go test file:

* Is any file within a package whose name ends with `_test.go`.
* Is executed when running the command:

```bash
go test ./...
```

### Go Test Functions

A Go test function:

* Must start with the prefix `Test`.
* Must take a single parameter of type `*testing.T`.

Example:

```go
func TestExample(t *testing.T) {
    // test logic here
}
```

### table-driven testing approach

We also adopt a **table-driven testing approach**, which makes tests more scalable, consistent, and easier to extend.

**Example01**
```go
package mathutils

import "testing"

// Function under test
func Add(a, b int) int {
    return a + b
}

// Table-driven test
func TestAdd(t *testing.T) {
    // Define test cases
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"with a negative number", -1, 4, 3},
        {"both negative numbers", -2, -3, -5},
    }

    // Run test cases
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

```
**Example02**
```go
func TestHelloHandler(t *testing.T) {
    tests := []struct {
        name       string
        url        string
        wantStatus int
        wantBody   string
    }{
        {"valid request", "/hello?name=Go", 200, "Hello, Go!"},
        {"missing name", "/hello", 400, "missing name"},
        {"not found", "/invalid", 404, "404 page not found"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", tt.url, nil)
            w := httptest.NewRecorder()

            router().ServeHTTP(w, req)

            if w.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d", w.Code, tt.wantStatus)
            }
            if !strings.Contains(w.Body.String(), tt.wantBody) {
                t.Errorf("body = %q; want %q", w.Body.String(), tt.wantBody)
            }
        })
    }
}
```



### What is a Test?

At its core, a test consists of three steps:

1. **Arrange** – provide **known inputs** and define the **expected outcome**.
2. **Act** – execute the **function or code** under test with the known inputs.
3. **Assert** – compare the **obtained result** with the **expected result**.

A test can:
- **pass**: if the **obtained result** and **expected result** match.
- **fail**: if they don’t, signaling a potential bug or unexpected behavior.

Depending on the context, the “obtained result” of a test might be:

- **A return value** (e.g., `Add(2,3)` should return `5`).
- **An error** (e.g., dividing by zero should return an error).
- **A side effect** (e.g., a file is created, a log message is written).
- **A performance constraint** (e.g., execution must complete under 100 ms).
- **A system state** (e.g., a flag is set, a resource is locked).
- **A performance metric** (e.g., memory usage or throughput).



## Commiting

Git commit messages, follows a [conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) style:

| Prefix      | Description                                      | Example                                       |
| ----------- | ------------------------------------------------ | --------------------------------------------- |
| `chore:`    | Maintenance or tooling work, no behavior change  | `chore: update CI workflow`                   |
| `feat:`     | Adds a new feature                               | `feat: add user login functionality`          |
| `fix:`      | Bug fix                                          | `fix: correct typo in API response`           |
| `docs:`     | Documentation changes                            | `docs: update README with setup instructions` |
| `refactor:` | Code changes that don’t add features or fix bugs | `refactor: simplify helper functions`         |
| `test:`     | Adding or updating tests                         | `test: add unit tests for file parser`        |

## TODO: Local Development vs. Shared Development environment
 

<!-- 

## Testing the code
As the Go community, we are using the **white-box testing** framework, allowing to test also private functions: 
- By **keeping** test files **alongside** the code. 
- It promote a **strong unit testing culture** 

we also use a table-driven approach by ...


A `GO` test file 
  - is simply any file in a **package** that ends with `_test.go`
  - is played whenever this cli occur : `go test ./...`

A `GO` test is a function that:
  - starts with the string `Test` 
  - takes `*testing.T` as an argument.  -->


----


# Release History & Changelog  

Track version updates and changes:  
- **📦 Latest Release**: `vX.X.X` ([GitHub Releases](#))  
- **📄 Full Changelog**: See [CHANGELOG.md](CHANGELOG.md) for detailed version history.  


