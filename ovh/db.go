package ovh

import (
	"path/filepath"
)

const DOMAIN_EU = "eu.api.ovh.com"
const NS_VERSION = "/v1"

var credentialRelPath = filepath.Join("wkspc", ".config", "ovh", "credential.json")
