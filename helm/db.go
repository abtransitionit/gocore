package helm

var MapHelmRepoReference = MapHelmRepo{
	"calico": {
		Url:  "https://docs.tigera.io/calico/charts",
		Name: "MapHelmRepo",
	},
	"bitnami": {
		Url:  "https://docs.tigera.io/calico/charts",
		Name: "MapHelmRepo",
		Doc: ["MapHelmRepo"],
	},
	"cilium": {
		Url:  "https://helm.cilium.io/",
		Name: "MapHelmRepo",
		Doc: ["MapHelmRepo"],
	},
}


// /Users/max/wkspc/git/lucg/cmd/kbe/helm

  // // Kubernetes/Helm CNI-Cilium conf
	// CiliumHelmRepoUrl   = "https://helm.cilium.io/"
	// CiliumHelmRepoName  = "cilium"
	// CiliumHelmChartName = "cilium"
	// CiliumHelmChartVersion = "1.17"  // compatible with K8s 1.32.x-1.33.x
	// CiliumK8sNamespace  = "kube-system"
  // CiliumHelmReleaseName = "kbe-cilium"
	// CiliumDocUrl        = "https://docs.cilium.io/"
	// CiliumGithub        = "https://github.com/cilium/cilium"

	// // Kubernetes/Helm Dashboard conf
	// KDashbHelmRepoUrl  = "https://kubernetes.github.io/dashboard/"
	// KDashbHelmRepoName = "kdashb"
	// KDashbHelmChartName = "kubernetes-dashboard"
	// KDashbHelmChartVersion = "7.12"  // works with K8s 1.32.x
	// KDashbK8sNamespace = "kdashb"
	// KDashbHelmReleaseName = "kbe-kdash"
	// KDashbDocUrl       = "https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/"
	// KDashbGithub       = "https://github.com/kubernetes/dashboard"

	// // Kubernetes/Helm Nginx controller conf
	// IngressNginxControllerHelmRepoUrl  = "https://kubernetes.github.io/ingress-nginx"
	// IngressNginxControllerHelmRepoName = "ingnginx"
	// IngressNginxControllerHelmChartName = "ingress-nginx"
	// IngressNginxControllerHelmChartVersion = "4.12"
	// IngressNginxControllerK8sNamespace = "ingress-nginx"
	// IngressNginxControllerHelmReleaseName = "kbe-ingress-nginx"
	// IngressNginxControllerDocUrl       = "https://github.com/kubernetes/ingress-nginx/blob/main/README.md#readme"
	// IngressNginxControllerDocUrl2      = "https://kubernetes.github.io/ingress-nginx/deploy/"
	// IngressNginxControllerGithub       = "https://github.com/kubernetes/ingress-nginx"
