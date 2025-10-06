package helm

type HelmRepo struct {
	Name string // logical name
	Desc string
	Url  string
	Doc  []string
}

type MapHelmRepo map[string]HelmRepo
