# Term
- SA = **S**ervice **A**ccount
- AT = **A**ccess **T**oken

# To know
- The service account is used by a client to request the `OVH` API
- The Service Account Id is the client Id
- The Service Account secret is the client Secret


# Example of file
`~/wkspc/.config/ovh/credential.json`

```json
{
  "myPolicy": [
    {
      "id": "ae5b786b-fe87-57b3-8b8e-0fa1cf1666a5"
    }
  ],
  "nic": {
    "id": "34eb8aa-ovh",
    "pwd": "enckshdè§te"
  },
  "serviceAccount": {
    "access_token": "xxxxxxx",
    "clientId": "EU.4332dff0faf4e0cc",
    "clientSecret": "4332dff0faf4e0cc4332dff0faf4e0cc",
    "description": "service account to be used by GO code to request OVH API",
    "flow": "CLIENT_CREDENTIALS",
    "identity": "urn:v1:eu:identity:credential:34eb8aa-ovh/oauth2-EU.4332dff0faf4e0cc",
    "name": "my_service_name"
  },
  "sshKeyId": "MySshKeyId"
}
```
`~/wkspc/.config/ovh/vps.json`
```json
{
	"vps01": {
		"DisplayName": "vm.ovh.01",
		"NameId":      "vps-4873ab3.vps.ovh.net",
		"Distro":      "ubuntu2504",
	},
	"vps02": {
		"DisplayName": "vm.ovh.02",
		"NameId":      "vps-aa12987.vps.ovh.net",
		"Distro":      "alma9",
	},
	"vps03": {
		"DisplayName": "vm.ovh.03",
		"NameId":      "vps-987af5e.vps.ovh.net",
		"Distro":      "rocky9",
	},
	"vps04": {
		"DisplayName": "vm.ovh.04",
		"NameId":      "vps-cab2764.vps.ovh.ca",
		"Distro":      "Fedora42",
	},
	"vps05": {
		"DisplayName": "vm.ovh.05",
		"NameId":      "vps-92afb534.vps.ovh.net",
		"Distro":      "debian12",
	}
}
```
