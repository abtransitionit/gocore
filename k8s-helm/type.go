package helm

type HelmRepo struct {
	Name string // logical name
	Desc string
	Url  string
	Doc  []string
}
type HelmChart struct {
	Name string
	Desc string
	Repo HelmRepo
}
type HelmRelease struct {
	Name          string
	Version       string
	Repo          HelmRepo
	ChartName     []string
	CharNamespace string
	ValueFile     string
}

type MapHelmRepo map[string]HelmRepo
