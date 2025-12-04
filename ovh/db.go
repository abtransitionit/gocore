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

var vpsOsImageReference = MapVpsOsImage{
	"alma9": {
		Name: "AlmaLinux 9",
		Id:   "4160a29f-bee5-41d9-865b-86acb6b03fe9",
	},
	"Fedora42": {Name: "Fedora 42",
		Id: "e3db5837-26ae-4a9e-bfab-344d5ad06388",
	},
	"ubuntu2504": {
		Name: "Ubuntu 25.04",
		Id:   "89581538-39cd-4855-9e37-f6b11d954d18",
	},
	"debian12": {
		Name: "Debian 12",
		Id:   "51fa6918-b1be-4922-a783-a989e1ff4925",
	},
	"rocky10": {
		Name: "Rocky Linux 10",
		Id:   "9449995a-1525-4bda-838e-d8c9d3f974d5",
	},
}

var endpointReference = apicli.MapEndpoint{
	"CreateToken":     {Verb: "POST", Desc: "ceate SA:token", Path: "/auth/oauth2/token"},
	"MeGetInfo":       {Verb: "GET", Desc: "get Me:Info", Path: fmt.Sprintf("/%s/me", NS_V1)},
	"ImageGetList":    {Verb: "GET", Desc: "get Image:List", Path: fmt.Sprintf("/%s/vps/{id}/images/available", NS_V1)},
	"ImageGetDetail":  {Verb: "GET", Desc: "get Image:List", Path: fmt.Sprintf("/%s/vps/{idv}/images/available/{idi}", NS_V1)},
	"VpsGetList":      {Verb: "GET", Desc: "get Vps:List", Path: fmt.Sprintf("/%s/vps", NS_V1)},
	"VpsGetDetail":    {Verb: "GET", Desc: "get Vps:Detail", Path: fmt.Sprintf("/%s/vps/{id}", NS_V1)},
	"VpsGetOs":        {Verb: "GET", Desc: "get Vps:Os", Path: fmt.Sprintf("/%s/vps/{id}/images/current", NS_V1)},
	"VpsReinstall":    {Verb: "POST", Desc: "re-install Vps", Path: fmt.Sprintf("/%s/vps/{id}/rebuild", NS_V1)},
	"SshKeyGetList":   {Verb: "GET", Desc: "get SSHKey:List", Path: fmt.Sprintf("/%s/me/sshKey", NS_V1)},
	"SshKeyGetDetail": {Verb: "GET", Desc: "get SSHKey:Detail", Path: fmt.Sprintf("/%s/me/sshKey/{id}", NS_V1)},
}

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
