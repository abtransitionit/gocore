package helm

var MapHelmRepoReference = MapHelmRepo{
	"calico": {
		Desc: "The calico operator Helm repository",
		Url:  "https://docs.tigera.io/calico/charts",
		Name: "calico",
	},
	"bitnami": {
		Desc: "The bitnami  Helm repository",
		Url:  "https://charts.bitnami.com/bitnami",
		Name: "bitnami",
	},
	"cilium": {
		Desc: "The cilium  Helm repository",
		Url:  "https://helm.cilium.io/",
		Name: "cilium",
		Doc:  []string{"https://github.com/cilium", "https://github.com/cilium/cilium-cli", "https://docs.cilium.io/en/stable/", "https://cilium.io/"},
	},
	"kdashb": {
		Desc: "The standard kubernetes dashboard Helm repository",
		Url:  "https://kubernetes.github.io/dashboard/",
		Name: "kdashb",
	},
	"ingressNginx": {
		Desc: "The nginx ingress controler Helm repository",
		Url:  "https://kubernetes.github.io/ingress-nginx",
		Name: "ingressNginx",
	},
}

// var ListHelmRepo = []HelmRepoProperty{
// 	{KDashbHelmRepoName, "The standard kubernetes dashboard Helm repository", KDashbHelmRepoUrl, new(bool)},
// 	{CiliumHelmRepoName, "The cilium   Helm repository", CiliumHelmRepoUrl, new(bool)},
// 	{"bitnami", "The bitnami  Helm repository", "https://charts.bitnami.com/bitnami", new(bool)},
// 	{IngressNginxControllerHelmRepoName, "The nginx ingress controler Helm repository", IngressNginxControllerHelmRepoUrl, new(bool)},
// }

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
