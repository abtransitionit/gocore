package ovh

import (
	"fmt"
	"path/filepath"

	"github.com/abtransitionit/gocore/apicli"
)

const DOMAIN_EU = "eu.api.ovh.com"
const DOMAIN_STD = "www.ovh.com"
const NS_V1 = "v1"
const NS_V2 = "v2"

var credentialRelPath = filepath.Join("wkspc", ".config", "ovh", "credential.json")
var listVpsRelPath = filepath.Join("wkspc", ".config", "ovh", "vps.json")

// var vpsReference = MapVpsInfo{
// 	"o1": {
// 		DisplayName: "vm.ovh.01",
// 		NameId:      "vps-9c33782a.vps.ovh.net",
// 		Distro:      "ubuntu2504",
// 	},
// 	"o2": {
// 		DisplayName: "vm.ovh.02",
// 		NameId:      "vps-ff9b9706.vps.ovh.net",
// 		Distro:      "alma9",
// 	},
// 	"o3": {
// 		DisplayName: "vm.ovh.03",
// 		NameId:      "vps-5b6ad1b7.vps.ovh.net",
// 		Distro:      "rocky9",
// 	},
// 	"o4": {
// 		DisplayName: "vm.ovh.04",
// 		NameId:      "vps-54f2f0c1.vps.ovh.ca",
// 		Distro:      "Fedora42",
// 	},
// 	"o5": {
// 		DisplayName: "vm.ovh.05",
// 		NameId:      "vps-a7a8f7f6.vps.ovh.net",
// 		Distro:      "debian12",
// 	},
// }

var vpsOsImageReference = MapVpsOsImage{
	"alma9": {
		Name: "AlmaLinux 9",
		Id:   "b8d58142-d568-4ff0-b5ef-10e9f1b65cf2",
	},
	"Fedora42": {Name: "Fedora 42",
		Id: "216873ec-b939-4b6b-9d6a-ec9d4c2aca33",
	},
	"ubuntu2504": {
		Name: "Ubuntu 25.04",
		Id:   "0f814e64-6ced-43d9-824b-467d46e00688",
	},
	"debian12": {
		Name: "Debian 12",
		Id:   "c8206049-3c17-4d6b-983a-3970f53b9ff5",
	},
	"rocky9": {
		Name: "Rocky Linux 9",
		Id:   "4862102e-029e-498e-917b-9fa16751c91d",
	},
}

var endpointReference = apicli.MapEndpoint{
	"CreateToken":     {Verb: "POST", Desc: "ceate SA:token", Path: "/auth/oauth2/token"},
	"MeGetInfo":       {Verb: "GET", Desc: "get Me:Info", Path: fmt.Sprintf("/%s/me", NS_V1)},
	"VpsGetList":      {Verb: "GET", Desc: "get Vps:List", Path: fmt.Sprintf("/%s/vps", NS_V1)},
	"VpsGetDetail":    {Verb: "GET", Desc: "get Vps:Detail", Path: fmt.Sprintf("/%s/vps/{id}", NS_V1)},
	"VpsGetOs":        {Verb: "GET", Desc: "get Vps:Os", Path: fmt.Sprintf("/%s/vps/{id}/images/current", NS_V1)},
	"VpsReinstall":    {Verb: "POST", Desc: "re-install Vps", Path: fmt.Sprintf("/%s/vps/{id}/rebuild", NS_V1)},
	"SshKeyGetList":   {Verb: "GET", Desc: "get SSHKey:List", Path: fmt.Sprintf("/%s/me/sshKey", NS_V1)},
	"SshKeyGetDetail": {Verb: "GET", Desc: "get SSHKey:Detail", Path: fmt.Sprintf("/%s/me/sshKey/{id}", NS_V1)},
}
