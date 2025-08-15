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


## Installation

To use this library in your project, run:

- `go` `get` [github.com/abtransitionit/gocore](https://github.com/abtransitionit/gocore)

---

# Getting started
## Howto modify interfaces
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
## Howto test the code
- For `GO` a test file is simply any file in a **package** that ends with `_test.go`
- For `GO` a way to play any **unit** test is : `go test ./...`

# Contributing  

We welcome contributions! Before participating, please review:  
- **[Code of Conduct](.github/CODE_OF_CONDUCT.md)** – Our community guidelines.  
- **[Contributing Guide](.github/CONTRIBUTING.md)** – How to submit issues, PRs, and more.  

## Local development env Vs. Shared development env
----


# Release History & Changelog  

Track version updates and changes:  
- **📦 Latest Release**: `vX.X.X` ([GitHub Releases](#))  
- **📄 Full Changelog**: See [CHANGELOG.md](CHANGELOG.md) for detailed version history.  

