package ovh

// Define structs for the request body and response
// The request body needs to be `x-www-form-urlencoded`
// so we'll use a url.Values map and convert it to a string.
type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	// TokenType   string `json:"token_type"`
}

type VpsOsImage struct {
	Name string
	Id   string
}
type MapVpsOsImage map[string]VpsOsImage

// type VpsInfo struct {
// 	DisplayName string
// 	NameId      string
// 	Distro      string
// }
// type MapVpsInfo map[string]VpsInfo

type CredentialStruct struct {
	SshKeyId       string `json:"sshKeyId"`
	ServiceAccount struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		AccessToken  string `json:"access_token"`
	} `json:"serviceAccount"`
}

type ListVpsStruct map[string]Vps

type VpsReinstallParam struct {
	DoNotSendPassword bool   `json:"doNotSendPassword"`
	ImageId           string `json:"imageId"`
	PublicSshKey      string `json:"publicSshKey"`
}

type ImageDetail struct {
	Name string
	Id   string
}

// -----------------------------------
// -- struct to contain YAML file ----
// -----------------------------------

// ##### VPS IMAGE DITRO INVENTORY #######

type DistroYaml struct {
	Distro []Distro
}

type Distro struct {
	Id   string
	Name string
}

// ##### VPS  INVENTORY #######

type VpsYaml struct {
	Vps map[string]Vps
}

// Description: represents a VPS
type Vps struct {
	Name        string
	NameId      string
	Distro      string
	NameDynamic string `json:"nameDynamic,omitempty"` // computed field
}

// type Vps struct {
// 	Name        string
// 	Id          string
// 	Distro      string
// 	NameDynamic string // computed field
// }
