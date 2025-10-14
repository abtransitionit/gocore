package cilium

// define the cilium configuration file
var configFileTpl = `
# Disable kube-proxy replacement for simplicity
kubeProxyReplacement: true # other choice: disable

k8sServiceHost: {{.K8sApiServerIp}}

k8sServicePort: 6443
# Enable Hubble (Cilium’s observability layer)
hubble:
  enabled: false      # You can set to true later if you want flow visibility

# configure Pods CIDR
ipam:
  mode: cluster-pool # other choice: kubernetes (Standard mode for most clusters)
  operator:
    clusterPoolIPv4PodCIDRList:
      - {{.K8sPodCidr}}
    clusterPoolIPv4MaskSize: 24
`
