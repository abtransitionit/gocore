package helm

type HelmRepo struct {
	Name string // logical name
	Url  string
	Doc  []string
}

type MapHelmRepo map[string]HelmRepo
