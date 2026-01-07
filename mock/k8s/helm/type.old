package helm

// define types
type Repo struct {
	Name string // eg. cilium, kdashb
	Desc string
	Url  string
	Doc  []string
}

// defrine slices
type RepoSlice []Repo

// define getters
func GetRepo(i Repo) *Repo {
	return &Repo{
		Name: i.Name,
		Url:  i.Url,
	}
}

// -------------------------------------------------------
// -------	 struct for YAML Repo List
// -------------------------------------------------------

// Description: represents the organization's repository db for helm Repo(s)
//
// Notes:
//   - Manage the YAML repo file
//   - denotes the whitelist
type MapYaml struct {
	List map[string]Repo
}

// type HelmRepo struct {
// 	Name string // logical name
// 	Desc string
// 	Url  string
// 	Doc  []string
// }
// type HelmChart struct {
// 	FullName string //ie. RepoName/ChartName or /tmp/chart/ChartName
// 	Version  string
// 	Desc     string
// 	Repo     HelmRepo
// }
// type HelmRelease struct {
// 	Name      string
// 	Repo      HelmRepo
// 	Chart     HelmChart
// 	Namespace string
// 	ValueFile string
// }

// type MapHelmRepo map[string]HelmRepo
