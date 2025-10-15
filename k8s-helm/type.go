package helm

type HelmRepo struct {
	Name string // logical name
	Desc string
	Url  string
	Doc  []string
}
type HelmChart struct {
	FullPath string //ie. /tmp/ChartName
	FullName string //ie. RepoName/ChartName
	Version  string
	Desc     string
	Repo     HelmRepo
}
type HelmRelease struct {
	Name      string
	Repo      HelmRepo
	Chart     HelmChart
	Namespace string
	ValueFile string
}

type MapHelmRepo map[string]HelmRepo
