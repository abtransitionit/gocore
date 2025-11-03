# Purpose
- return a struct that is a merge of several `YAML` configuration file (`conf.yaml`)
- the YMAL config file can be found at different location (package + Global + working dir)

# `Conf.yaml` 
- can be common to a set of `workflow`s
- can be specific to one `workflow`


example of a config `YAML` file specicif to one workflow named `kbe`:

```yaml
workflow:
  kbe:
    customRcFileName: ".profile.luc"
    binFolderPath: "/usr/local/bin"

    node:
      all: ["o1u", "o2a", "o3r", "o4f", "o5d"]
      worker: ["o2a", "o3r", "o4f", "o5d"]

    goCli:
      controlPlane:
      - { name: "helm", version: "3.17.3" }
        
    da:
      repo:
        node:
          - { name: "k8s", fileName: "kbe-k8s", version: "1.32" }
      pkg:
        node:
          - { name: "kubelet" }
        controlPlane:  
          - { name: "kubectl" }
        required:
          - { name: "gnupg" }

    helm:
      repo:
        - { name: "kdabsh" }
      release:
        - name: "kbe-cilium"
          chart: cilium
          repo:  cilium
          version: "0.18.7"
          namespace: "kube-system"

    cluster:
      podCidr: "192.168.0.0/16"
      serviceCidr: "172.16.0.0/16"
      crSocketName: "unix:///var/run/crio/crio.sock"
      k8sVersion:   "1.32.0"

```

example of a config `YAML` common to a set of workflows name kbe:

```yaml
workflow:
  kbe:
    customRcFileName: ".profile.luc"
    binFolderPath: "/usr/local/bin"
  kind:
    node:
      all: ["o1u", "o2a", "o3r", "o4f", "o5d"]
      worker: ["o2a", "o3r", "o4f", "o5d"]
```      

# Todo
```
config.GetStringSlice("node.all")
config.Get("da.repo.node")
config.Get("da.pkg.node")
config.Get("helm.release")
config.Get("cluster")
```