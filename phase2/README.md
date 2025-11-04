# Purpose
- A framework to provision (ie. execute actions on) `VMs` (**V**irtual **M**acchine**s**) or locally
- Actions are `GO` functions that may 
  - have dependencies
  - run concurently or locally
  - are described/defined in a YAML configuration file
- The whole action is named a c and is described/defined in a YAML configuration file

# How it works
- a `GO` `cobra` `CLI` allows to list, print, run a part or the whole of any registred workflows
- the workflows are hard coded in the CLI
- the workflow's default configuration is hard coded in the CLI **and** can be overriden

# Todo
- Extend the concept of **workflow** to containers that may be local or remote.
- use the CLI as a CI/CD pipeline