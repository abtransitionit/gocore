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



# Process

```go
	var err error
	wkf, err = corephase.NewWorkflowFromPhases(
		corephase.NewPhase("checkVmAccess", "Check if VMs are SSH reachable", vm.CheckVmSshAccess, nil),
		corephase.NewPhase("copyAgent", "copy LUC CLI agent to all VMs", luc.DeployLuc, []string{"checkVmAccess"}),
		corephase.NewPhase("upgradeOs", "provision OS nodes with latest dnfapt packages and repositories.", dnfapt.UpgradeVmOs, []string{"copyAgent"}),
		corephase.NewPhase("updateApp", "provision required/missing standard dnfapt packages.", dnfapt.UpdateVmOsApp(listRequiredDaPackage), []string{"upgradeOs"}),
		corephase.NewPhase("installDaRepository", "provision Dnfapt package repositor(y)(ies).", dnfapt.InstallDaRepository(cfg.Da.Repo.Node), []string{"updateApp"}),
		corephase.NewPhase("installDaPackage", "provision Dnfapt package(s) on all nodes.", dnfapt.InstallDaPackage(cfg.Da.Pkg.Node), []string{"installDaRepository"}),
		corephase.NewPhase("installDaPackageCplane", "provision Dnfapt package(s) on CPlane only.", dnfapt.InstallDaPackage(cfg.Da.Pkg.ControlPlane, targetsCP), []string{"installDaPackage"}),
		corephase.NewPhase("loadOsKernelModule", "load OS kernel module(s).", taskoskernel.LoadOsKModule(sliceOsKModule, kernelFilename), []string{"installDaPackage"}),
		corephase.NewPhase("loadOsKernelParam", "set OS kernel paramleter(s).", taskoskernel.LoadOsKParam(sliceOsKParam, kernelFilename), []string{"loadOsKernelModule"}),
		corephase.NewPhase("confSelinux", "Configure Selinux.", selinux.ConfigureSelinux(), []string{"loadOsKernelParam"}),
		corephase.NewPhase("enableOsService", "enable OS services to start after a reboot", oservice.EnableOsService(sliceOsServiceEnable), []string{"confSelinux"}),
		corephase.NewPhase("startOsService", "start OS services for current session", oservice.StartOsService(sliceOsServiceStart), []string{"confSelinux"}),
		corephase.NewPhase("resetCPlane", "reset the control plane(s).", taskk8s.ResetNode(targetsCP), []string{"startOsService"}),
		corephase.NewPhase("initCPlane", "initialize the control plane(s) (aka. boostrap the cluster).", taskk8s.InitCPlane(targetsCP, cfg.Cluster), []string{"resetCPlane"}),
		corephase.NewPhase("resetWorker", "reset the workers(s).", taskk8s.ResetNode(targetsWorker), []string{"initCPlane"}),
		corephase.NewPhase("addWorker", "Add the K8s worker(s) to the K8s cluster.", taskk8s.AddWorker(targetsCP[0], targetsWorker), []string{"resetWorker"}),
		corephase.NewPhase("confKubectlOnCPlane", "Configure kubectl on the control plane(s).", taskk8s.ConfigureKubectlOnCplane(targetsCP[0]), []string{"resetWorker"}),
		corephase.NewPhase("installGoCliCplane", "provision Go CLI(s).", gocli.InstallGoCliOnVm(listGoCli, targetsCP), []string{"confKubectlOnCPlane"}),
		corephase.NewPhase("createRcFile", "create a custom RC file in user's home.", util.CreateCustomRcFile(customRcFileName), []string{"installGoCliCplane"}),
		corephase.NewPhase("setPathEnvar", "configure PATH envvar into current user's custom RC file.", util.SetPath(binFolderPath, customRcFileName), []string{"createRcFile"}),
		corephase.NewPhase("installHelmRepo", "install Helm chart repositories.", helm.InstallHelmRepo(sliceHelmRepo), []string{"setPathEnvar"}),
		corephase.NewPhase("createCiliumRelease", "install and configure the CNI: Cilium on all nodes.", task_cilium.InstallCniCilium, []string{"installHelmRepo"}),
	)
```



# `Execute()` is the orchestrator engine

Even though it only logs today, its *intended role* is clearly to:

* Build the dependency graph
* Compute tiers
* Run phases concurrently within tiers
* Inject config variables as function arguments
* Dispatch functions onto nodes
* Handle errors, retries, and maybe rollback

Todo

* Logs workflow name + description
* Logs the rule that phases in the same tier run concurrently
* Returns `nil` (success)
* Does *not*:

  * resolve nodes
  * resolve parameters
  * compute tiers
  * run phases
  * run functions
  * check dependencies
  * perform SSH
  * handle concurrency
