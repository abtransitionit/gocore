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
var vpsListRelPath = filepath.Join("wkspc", ".config", "ovh", "vps.json")

var endpointReference = apicli.MapEndpoint{
	"CreateToken":     {Verb: "POST", Desc: "ceate SA:token", Path: "/auth/oauth2/token"},
	"MeGetInfo":       {Verb: "GET", Desc: "get Me:Info", Path: fmt.Sprintf("/%s/me", NS_V1)},
	"VpsImageGetList": {Verb: "GET", Desc: "get Image:List", Path: fmt.Sprintf("/%s/vps/{id}/images/available", NS_V1)},
	"ImageGetDetail":  {Verb: "GET", Desc: "get Image:List", Path: fmt.Sprintf("/%s/vps/{idv}/images/available/{idi}", NS_V1)},
	"VpsGetList":      {Verb: "GET", Desc: "get Vps:List", Path: fmt.Sprintf("/%s/vps", NS_V1)},
	"VpsGetDetail":    {Verb: "GET", Desc: "get Vps:Detail", Path: fmt.Sprintf("/%s/vps/{id}", NS_V1)},
	"VpsGetOs":        {Verb: "GET", Desc: "get Vps:Os", Path: fmt.Sprintf("/%s/vps/{id}/images/current", NS_V1)},
	"VpsReinstall":    {Verb: "POST", Desc: "re-install Vps", Path: fmt.Sprintf("/%s/vps/{id}/rebuild", NS_V1)},
	"SshKeyGetList":   {Verb: "GET", Desc: "get SSHKey:List", Path: fmt.Sprintf("/%s/me/sshKey", NS_V1)},
	"SshKeyGetDetail": {Verb: "GET", Desc: "get SSHKey:Detail", Path: fmt.Sprintf("/%s/me/sshKey/{id}", NS_V1)},
}

var vpsOsImageReference = MapVpsOsImage{
	"alma9": {
		Name: "AlmaLinux 9",
		Id:   "4160a29f-bee5-41d9-865b-86acb6b03fe9",
	},
	"debian13": {
		Name: "Debian 13",
		Id:   "24195a3e-7b12-44ef-ade3-22bcf245e7a7",
	},
	"Fedora42": {Name: "Fedora 42",
		Id: "e3db5837-26ae-4a9e-bfab-344d5ad06388",
	},
	"rocky9": {
		Name: "Rocky Linux 10",
		Id:   "ccb9c787-3f0e-4590-a507-7a50b7261c8c",
	},
	"ubuntu2504": {
		Name: "Ubuntu 25.04",
		Id:   "89581538-39cd-4855-9e37-f6b11d954d18",
	},
}
